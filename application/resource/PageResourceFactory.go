package resource

type PageResourceFactory struct {
	LoginUrl  string
	ImageUrl  string
	DesignUrl string
}

func NewPageResourceFactory(loginUrl string, designUrl string, imageUrl string) *PageResourceFactory {

	return &PageResourceFactory{
		LoginUrl:  loginUrl,
		DesignUrl: designUrl,
		ImageUrl:  imageUrl}
}

func (s *PageResourceFactory) CreateResource() *PageResource {
	pageResource := &PageResource{}
	pageResource.ImageUrl = s.ImageUrl
	pageResource.DesignUrl = s.DesignUrl
	pageResource.LoginUrl = s.LoginUrl
	pageResource.ErrorMessages = make([]string, 0)
	return pageResource
}
