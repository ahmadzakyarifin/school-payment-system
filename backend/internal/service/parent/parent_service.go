package parentservice

type ParentService interface {
	// Fitur spesifik wali murid nanti di sini
}

type parentService struct {
}

func New() ParentService {
	return &parentService{}
}
