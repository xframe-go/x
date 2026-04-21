package repository

type PrimaryKeyGetter[M any, K comparable] func(m M) K
