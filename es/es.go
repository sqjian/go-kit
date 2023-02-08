package es

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"os"
)

type Cli struct {
	meta struct {
		cli   *http.Client
		debug bool
		hosts []string
	}

	cli *elastic.Client

	debugLogger elastic.Logger
}

func newEsCli(opts ...OptionFunc) (*Cli, error) {
	cli := &Cli{}

	for _, opt := range opts {
		opt(cli)
	}

	if len(cli.meta.hosts) == 0 {
		return nil, ErrWrapper(IllegalParams)
	}
	if cli.meta.debug {
		cli.debugLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	esCli, err := elastic.NewClient(
		elastic.SetURL(cli.meta.hosts...),
		elastic.SetSniff(false),
		elastic.SetTraceLog(cli.debugLogger),
		elastic.SetHttpClient(cli.meta.cli),
	)
	if err != nil {
		return nil, err
	}

	cli.cli = esCli

	return cli, nil
}

func (es *Cli) indexExists(ctx context.Context, index string) (bool, error) {
	return es.
		cli.
		IndexExists(index).
		Do(ctx)
}

func (es *Cli) createIndex(ctx context.Context, index string) error {
	_, rstErr := es.
		cli.
		CreateIndex(index).
		Do(ctx)

	return rstErr
}

func (es *Cli) writeDocs(ctx context.Context, index string, docs ...map[string]interface{}) error {
	switch {
	case len(docs) == 1:
		{
			_, err := es.
				cli.Index().
				Index(index).
				BodyJson(docs[0]).
				Do(ctx)
			return err
		}
	case len(docs) > 1:
		{
			bulkRequest := es.cli.Bulk()
			for _, doc := range docs {
				bulkRequest = bulkRequest.
					Add(
						elastic.
							NewBulkIndexRequest().
							Index(index).
							Doc(doc),
					)
			}
			_, err := bulkRequest.Do(ctx)
			return err
		}
	default:
		{
			return ErrWrapper(IllegalParams)
		}
	}

}

func (es *Cli) updateDocs(ctx context.Context, index string, docs ...struct {
	id  string
	doc map[string]interface{}
}) error {
	switch {
	case len(docs) == 1:
		{
			_, err := es.
				cli.
				Update().
				Index(index).
				Id(docs[0].id).
				Doc(docs[0].doc).
				Do(context.Background())
			return err
		}
	case len(docs) > 1:
		{
			bulkRequest := es.cli.Bulk()
			for _, doc := range docs {
				bulkRequest = bulkRequest.
					Add(
						elastic.
							NewBulkUpdateRequest().
							Index(index).
							Id(doc.id).
							Doc(doc.doc),
					)
			}
			_, err := bulkRequest.Do(ctx)
			return err
		}
	default:
		{
			return ErrWrapper(IllegalParams)
		}
	}
}

func (es *Cli) queryDocs(ctx context.Context, index string, condition map[string]interface{}) ([]string, error) {

	boolQuery := elastic.NewBoolQuery()
	for k, v := range condition {
		boolQuery = boolQuery.Filter(
			elastic.NewTermQuery(k, v),
		)
	}

	rst, rstErr := es.
		cli.
		Search().
		Index(index).
		Query(boolQuery).
		Do(ctx)

	if rstErr != nil {
		return nil, rstErr
	}

	var data []string
	if rst.TotalHits() > 0 {
		for _, doc := range rst.Hits.Hits {
			data = append(data, string(doc.Source))
		}
	}
	return data, nil
}
