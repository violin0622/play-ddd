package ptr

func To[A any](a A) *A { return &a }
