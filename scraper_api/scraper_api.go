package scraper_api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	// "golang.org/x/sync/errgroup"
)

// LinkedIn Credentials
var (
	email    = "kiokovincent12@gmail.com"
	password = "Kioko@1024"
)

// ScrapedData stores the scraped data
type ScrapedData struct {
	UserID string     `json:"user_id"`
	Data   []PostData `json:"data"`
}

// PostData stores individual post data
type PostData struct {
	PostID     string      `json:"post_id"`
	Likes      int         `json:"likes"`
	LikersData []LikerData `json:"likers_data"`
}

// LikerData stores data of individual likers
type LikerData struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
}

// var loggedInCtx context.Context // Cache authenticated context (optional)

func getLinkedInData(userID string) (ScrapedData, error) {
	var scrapedData ScrapedData

	loginURL := "https://www.linkedin.com/login?fromSignIn=true&trk=guest_homepage-basic_nav-header-signin"
	profileURL := fmt.Sprintf("https://www.linkedin.com/in/%s/", userID)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var err error
	// var html string

	log.Println("Logging in...")
	// Perform login
	err = chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.Sleep(2000*time.Millisecond),
		// extract the raw HTML from the page
		chromedp.SendKeys(`#username`, email, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Submit(`#password`, chromedp.ByID),
		chromedp.WaitVisible(`#voyager-feed`, chromedp.ByID),
	)
	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}

	// loggedInCtx = ctx // Cache context for future calls

	log.Println("Fetching profile...")
	log.Println(profileURL)

	// Navigate to profile page
	var finalActivityUrl string
	err = chromedp.Run(ctx,
		chromedp.Navigate(profileURL),
		chromedp.WaitVisible(`#global-nav`, chromedp.ByID),
		chromedp.AttributeValue(`footer a[data-test-app-aware-link]`, "href", &finalActivityUrl, nil, chromedp.ByQuery),
	)

	log.Println(finalActivityUrl)

	err = chromedp.Run(ctx,
		chromedp.Navigate(finalActivityUrl),
		chromedp.WaitVisible(`#global-nav`, chromedp.ByID),
	)

	if err != nil {
		log.Println("error: ", err)
		return scrapedData, fmt.Errorf("failed to navigate to profile: %w", err)
	}
	// now navigate to all activity

	log.Println("Loading and scrolling...")
	// Scroll to load more posts
	for i := 0; i < 5; i++ {
		chromedp.Run(ctx,
			chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil),
			chromedp.Sleep(1000*time.Millisecond),
		)
	}

	log.Println("Following to containers")

	// Scrape post containers
	var containers []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// select the root node on the page
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			// html, _ := dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			// fmt.Println(html)
			cont_nodes, err := dom.QuerySelectorAll(rootNode.NodeID, ".feed-shared-update-v2").Do(ctx)
			for _, nodeID := range cont_nodes {
				node, _ := dom.DescribeNode().WithNodeID(nodeID).Do(ctx)
				containers = append(containers, node)
			}
			return err
		}),
	)
	if err != nil {
		return scrapedData, fmt.Errorf("failed to get post containers: %w", err)
	}

	log.Println("Creating scrapping instances for concurrency... ")
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(containers))

	for _, container := range containers {
		wg.Add(1)
		go func(container *cdp.Node) {
			defer wg.Done()
			if err := scrapePost(ctx, container, &scrapedData, &mu); err != nil {
				errChan <- err
			}
		}(container)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return scrapedData, err
		}
	}

	scrapedData.UserID = userID
	return scrapedData, nil
}

func scrapePost(ctx context.Context, container *cdp.Node, scrapedData *ScrapedData, mu *sync.Mutex) error {
	log.Println("Scraping instance...")
	var postID string
	for _, attr := range container.Attributes {
		if strings.Contains(attr, "activity") {
			postID = attr
			break
		}
	}
	if postID == "" {
		return nil
	}

	var likersData []LikerData

	// Navigate to the specific LinkedIn post
	postReactionsURL := fmt.Sprintf("https://www.linkedin.com/analytics/post/%s/?resultType=REACTIONS", postID)
	err := chromedp.Run(ctx,
		chromedp.Navigate(postReactionsURL),
		chromedp.WaitVisible(`#global-nav`, chromedp.ByID),
	)
	if err != nil {
		return err
	}

	// Extract likers' data
	var likers []*cdp.Node
	err = chromedp.Run(ctx, chromedp.Nodes(`ul[aria-label="People who reacted"] li`, &likers))
	if err != nil {
		return err
	}

	for _, liker := range likers {
		var name, userID, title string
		chromedp.Run(ctx,
			chromedp.Text(`span[dir="ltr"] span`, &name, chromedp.ByQuery, chromedp.FromNode(liker)),
			chromedp.AttributeValue(`a`, "href", &userID, nil, chromedp.ByQuery, chromedp.FromNode(liker)),
			chromedp.Text(`div.artdeco-entity-lockup__subtitle`, &title, chromedp.ByQuery, chromedp.FromNode(liker)),
		)

		// Clean up userID
		standardB64Str := strings.Split(userID, "/")[4]
		//decoded, _ := base64.StdEncoding.DecodeString(standardB64Str)

		likersData = append(likersData, LikerData{
			Name:   name,
			UserID: standardB64Str,
			Title:  title,
		})
	}

	mu.Lock()
	scrapedData.Data = append(scrapedData.Data, PostData{
		PostID:     postID,
		Likes:      0,
		LikersData: likersData,
	})
	mu.Unlock()

	return nil
}

// Handler to get all users
func ScrapeDataView(c *gin.Context) {
	info, _ := getLinkedInData("ACoAACby69gBA97b1q6Dovl2Rb5F8pvE1N8hB28")
	c.JSON(http.StatusOK, gin.H{"details": info})
}
