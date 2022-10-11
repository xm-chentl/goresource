package postgres

// func Test_Connections(t *testing.T) {
// 	connStr := "postgres://test:123456@192.168.0.200:5432/didagogo"
// 	// connectConfig := pgx.ConnConfig{}
// 	config, err := pgxpool.ParseConfig(connStr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	config.MaxConns = 800
// 	config.MinConns = 10
// 	pool, err := pgxpool.ConnectConfig(context.Background(), config)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for i := 0; i < 1000; i++ {
// 		go func() {
// 			conn, err := pool.Acquire(context.Background())
// 			if err != nil {
// 				return
// 			}
// 			time.Sleep(1 * time.Second)
// 			conn.Release()
// 		}()
// 	}
// 	time.Sleep(30 * time.Second)
// }
