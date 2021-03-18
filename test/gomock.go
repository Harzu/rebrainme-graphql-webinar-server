package test

//go:generate mockgen -source=./../internal/repositories/orders/repository.go -destination=./mocks/ordersRepoMocks/orders_repo_mock.go -package=ordersRepoMocks
//go:generate mockgen -source=./../internal/repositories/users/repository.go -destination=./mocks/usersRepoMocks/users_repo_mock.go -package=usersRepoMocks
//go:generate mockgen -source=./../internal/repositories/products/repository.go -destination=./mocks/productsRepoMocks/products_repo_mock.go -package=productsRepoMocks

//go:generate mockgen -source=./../internal/services/storage/sessions.go -destination=./mocks/sessionStorageMocks/session_storage_mock.go -package=sessionStorageMocks
