package api

import (
	"context"
	"errors"

	fs "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"server/common/constant"
)

type queryopts struct {
	Collection string
	ID         string
	Cursor     int
	Category   int
	Limit      int
}

// getFeeds returns a list of feeds
func getFeeds(ctx context.Context, client *fs.Client) []map[string]interface{} {
	iter := client.Collection(constant.Feeds).Documents(ctx)
	var items []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		data := doc.Data()
		items = append(items, map[string]interface{}{
			"id":    data["id"],
			"title": data["title"],
		})
	}
	return items
}

// getEntry returns single entry based on specified collection name and entry ID
func getEntry(ctx context.Context, client *fs.Client, opts queryopts) (map[string]interface{}, error) {
	if opts.Collection == "" || opts.ID == "" {
		return nil, errors.New("missing Collection or ID in queryopts")
	}

	entry, err := client.Collection(opts.Collection).Doc(opts.ID).Get(ctx)
	if err != nil {
		return nil, err
	}
	return entry.Data(), nil
}

func getEntries(ctx context.Context, client *fs.Client, opts queryopts) []map[string]interface{} {
	// default limit 10
	// max limit 20
	if opts.Limit == 0 {
		opts.Limit = 10
	} else if opts.Limit > 20 {
		opts.Limit = 20
	}

	query := client.Collection(opts.Collection).OrderBy("published_at", fs.Desc).Limit(opts.Limit)
	if opts.Cursor > 0 {
		query = query.StartAfter(opts.Cursor)
	}
	if opts.Category > 0 {
		query = query.Where("category_id", "==", opts.Category)
	}

	iter := query.Documents(ctx)
	var items []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}
		items = append(items, doc.Data())
	}
	return items
}
