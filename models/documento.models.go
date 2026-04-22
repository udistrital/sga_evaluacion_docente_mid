package models

type MetaFile struct {
	Name            string  `json:"name"`
	MimeType        string  `json:"mime-type"`
	Encoding        *string `json:"encoding"`
	DigestAlgorithm string  `json:"digestAlgorithm"`
	Digest          string  `json:"digest"`
	Length          string  `json:"length"`
	Data            string  `json:"data"`
}

type Document struct {
	UIDUID               *string       `json:"uid:uid"`
	UIDMajorVersion      int           `json:"uid:major_version"`
	UIDMinorVersion      int           `json:"uid:minor_version"`
	ThumbThumbnail       MetaFile      `json:"thumb:thumbnail"`
	FileContent          MetaFile      `json:"file:content"`
	CommonIconExpanded   *string       `json:"common:icon-expanded"`
	CommonIcon           string        `json:"common:icon"`
	FilesFiles           []interface{} `json:"files:files"`
	DCDescription        *string       `json:"dc:description"`
	DCLanguage           *string       `json:"dc:language"`
	DCCoverage           *string       `json:"dc:coverage"`
	DCValid              *string       `json:"dc:valid"`
	DCCreator            string        `json:"dc:creator"`
	DCModified           string        `json:"dc:modified"`
	DCLastContributor    string        `json:"dc:lastContributor"`
	DCRights             *string       `json:"dc:rights"`
	DCExpired            *string       `json:"dc:expired"`
	DCFormat             *string       `json:"dc:format"`
	DCCreated            string        `json:"dc:created"`
	DCTitle              string        `json:"dc:title"`
	DCIssued             *string       `json:"dc:issued"`
	DCNature             *string       `json:"dc:nature"`
	DCSubjects           []string      `json:"dc:subjects"`
	DCContributors       []string      `json:"dc:contributors"`
	DCSource             *string       `json:"dc:source"`
	DCPublisher          *string       `json:"dc:publisher"`
	RelatedTextResources []interface{} `json:"relatedtext:relatedtextresources"`
	NXTagTags            []interface{} `json:"nxtag:tags"`
	File                 string        `json:"file"`
}
