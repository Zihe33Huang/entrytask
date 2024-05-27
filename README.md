# Project Structure




|── backend
│   ├── communication // protubuf
│   ├── http // http server
│   ├── rpc  // simple rpc framework
│   └── tcp  // tcp server
├── frontend
├── test    // locust for performance test




# Front End

1、 Navigate to front end project directory.  cd fronted .

2、 Install Dependencies.   npm install .

3、 Start project using  npm start .




# Back End

1、 Navigate to tcp server project main directory.  cd backend/tcp/main .

2、Start TCP server. go run main.go .

3、Navigate to http server project main directory.  cd backend/http/main .

4、Start HTTP server. go run main.go .

PS:  HTTP server must start after TCP server.
