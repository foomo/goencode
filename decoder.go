package goencode

// Decoder decodes target T back into source S.
type Decoder[S, T any] func(t T, s *S) error
