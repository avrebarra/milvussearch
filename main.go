package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/milvus-io/milvus-sdk-go/milvus"
	"github.com/urfave/cli"
)

func main() {
	cmdroot := &cli.App{
		Name:    "milvussearch",
		Version: "v1",
		Commands: []cli.Command{
			{
				Name: "checkcollection",
				Action: func(c *cli.Context) error {
					return CheckCollection()
				},
			},
			{
				Name: "setupcollection",
				Action: func(c *cli.Context) error {
					return SetupCollection()
				},
			},
			{
				Name: "setupcontent",
				Action: func(c *cli.Context) error {
					return SetupContents()
				},
			},
			{
				Name: "insert",
				Action: func(c *cli.Context) error {
					return PerformInsert()
				},
			},
			{
				Name: "search",
				Action: func(c *cli.Context) error {
					return PerformSearch()
				},
			},
		},
	}

	// execute
	if err := cmdroot.Run(os.Args); err != nil {
		fmt.Println("unexpected error:", err.Error())
		return
	}
}

var (
	client          milvus.MilvusClient
	hostaddr        = "178.128.97.181"
	hostport        = "19530"
	collectionname  = "test_experiment_1"
	vectordimension = 100
)

func init() {
	var err error

	log.SetFlags(0)

	log.Println("setting up connection")
	clientgrpc := milvus.Milvusclient{}
	client = milvus.NewMilvusClient(clientgrpc.Instance)
	err = client.Connect(milvus.ConnectParam{
		IPAddress: hostaddr,
		Port:      hostport,
	})
	if err != nil || !client.IsConnected() {
		err = fmt.Errorf("cannot connect to milvus: %w", err)
		panic(err)
	}
}

func CheckCollection() (err error) {
	// check collection
	log.Println("checking existing collection:", collectionname)
	hasCollection, status, err := client.HasCollection(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot check collection: %w", err)
		return
	}

	if !hasCollection {
		log.Println("collection not exist")
		return
	}

	// describe collection stats
	stat, status, err := client.GetCollectionStats(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot get collection stats: %w", err)
		return
	}
	log.Println(fmt.Sprintf("current collection stat: %s", stat))

	// describe index
	var inf milvus.IndexParam
	inf, status, err = client.GetIndexInfo(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot create index: %w", err)
		return
	}
	log.Println("- name:", inf.CollectionName)
	log.Println("- type:", inf.IndexType)
	log.Println("- extra_params:", inf.ExtraParams)

	// check collection entities count
	count, status, err := client.CountEntities(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot get collection info: %w", err)
		return
	}
	log.Println(fmt.Sprintf("current collection entities count: %d", count))

	return
}

func SetupCollection() (err error) {
	// check collection
	log.Println("checking existing collection:", collectionname)
	hasCollection, status, err := client.HasCollection(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot check collection: %w", err)
		return
	}

	if !hasCollection {
		// create collection
		log.Println("collection not exist: creating")
		log.Println("creating collection")
		status, err = client.CreateCollection(milvus.CollectionParam{
			CollectionName: collectionname,
			Dimension:      int64(vectordimension),
			IndexFileSize:  100,
			MetricType:     int64(milvus.L2),
		})
		if err != nil || !status.Ok() {
			err = fmt.Errorf("cannot create collection: %w", err)
			return
		}

		// create index
		log.Println("creating collection index")
		status, err = client.CreateIndex(&milvus.IndexParam{
			CollectionName: collectionname,
			IndexType:      milvus.IVFFLAT,
			ExtraParams:    `{"nlist" : 16384}`,
		})
		if err != nil || !status.Ok() {
			err = fmt.Errorf("cannot create index: %w", err)
			return
		}
	}

	// describe index
	var inf milvus.IndexParam
	inf, status, err = client.GetIndexInfo(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot create index: %w", err)
		return
	}
	log.Println("- name:", inf.CollectionName)
	log.Println("- type:", inf.IndexType)
	log.Println("- extra_params:", inf.ExtraParams)

	// preload collection
	status, err = client.LoadCollection(milvus.LoadCollectionParam{
		CollectionName:   collectionname,
		PartitionTagList: nil,
	})
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot load collection: %w", err)
		return
	}

	return
}

func SetupContents() (err error) {
	// insert random data until reach 1mil (maximum)
	countperbatch := 100

	// check collection stats
	count, status, err := client.CountEntities(collectionname)
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot get collection info: %w", err)
		return
	}

	iterlimit := 1000000 - count

	log.Println(fmt.Sprintf("current collection entities count: %d", count))
	log.Println(fmt.Sprintf("will create bloat %d entities", iterlimit))

	for {
		if iterlimit <= 0 {
			break
		}

		log.Println(fmt.Sprintf("inserting vectors to collection: %s", collectionname))
		log.Println(fmt.Sprintf("creating bloat vectors"))
		records := make([]milvus.Entity, countperbatch)
		for i := range records {
			data := make([]float32, vectordimension)
			for i := range data {
				data[i] = rand.Float32()
			}
			records[i] = milvus.Entity{FloatData: data}
		}
		log.Println(fmt.Sprintf("done creating bloat vectors"))

		log.Println(fmt.Sprintf("inserting bloat vectors: %d", countperbatch))
		starttime := time.Now()
		_, status, err := client.Insert(&milvus.InsertParam{
			CollectionName: collectionname,
			RecordArray:    records,
		})

		log.Println(fmt.Sprintf("done inserting bloat %d vectors (need more %d): took %s", countperbatch, iterlimit, time.Now().Sub(starttime).String()))
		if err != nil || !status.Ok() {
			err = fmt.Errorf("cannot insert data to collection: %w", err)
			log.Println(err)
		}

		iterlimit--
	}

	return
}

func PerformInsert() (err error) {
	// read entity from insertfile
	e := Entity{}
	err = e.FromFile("f_insert.yaml")
	if err != nil {
		return
	}

	// create milvus entity
	me := milvus.Entity{FloatData: make([]float32, 100)}
	for i, v := range e.Traits {
		me.FloatData[i] = v.Tendency
	}

	// insert
	ids, status, err := client.Insert(&milvus.InsertParam{
		CollectionName: collectionname,
		RecordArray:    []milvus.Entity{me},
	})
	if err != nil {
		err = fmt.Errorf("cannot insert data to collection: %w", err)
		return
	}
	if !status.Ok() {
		err = fmt.Errorf("cannot insert data to collection: %s", status.GetMessage())
		return
	}

	fmt.Println(fmt.Sprintf("INSERTED! ID %s got milvus id %d (hash:%s)", e.ID, ids[0], SHA1(fmt.Sprint(ids[0]))[0:5]))

	return
}

func PerformSearch() (err error) {
	// read entity from searchfile
	e := Entity{}
	err = e.FromFile("f_search.yaml")
	if err != nil {
		return
	}

	// create search data
	searchvector := make([]float32, vectordimension)
	for i := range e.Traits {
		searchvector[i] = e.Traits[i].Tendency
	}

	// search
	log.Println("searching collection")
	timestampsearch := time.Now()
	results, status, err := client.Search(milvus.SearchParam{
		CollectionName: collectionname,
		QueryEntities: []milvus.Entity{
			{FloatData: searchvector},
			// {FloatData: searchvector},
		},
		Topk:        10,
		ExtraParams: `{"nprobe" : 32}`,
	})

	log.Println(fmt.Sprintf("search done (took %s)", time.Now().Sub(timestampsearch).String()))
	if err != nil || !status.Ok() {
		err = fmt.Errorf("cannot search collection: %w", err)
		return
	}

	log.Println("search results: ")
	for _, r := range results.QueryResultList {
		for j, id := range r.Ids {
			log.Println(fmt.Sprintf("- result no %d: ID:%d (hash:%s)\tDIST:%f", j, id, SHA1(fmt.Sprint(id))[0:5], r.Distances[j]))
		}
	}

	return
}

func catcherr(err error) {
	if err != nil {
		err = fmt.Errorf("unexpected error: %w", err)
		log.Fatal(err)
	}
}
