package root

type Crypto interface {
  Salt(s string) (error, string)
  Compare(hash string, s string) (error, bool)
}