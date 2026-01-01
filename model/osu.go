package model

type BeatmapSearchResponse struct {
	BeatmapSets           []BeatmapSet      `json:"beatmapsets"`
	SearchInfo            map[string]string `json:"search"`
	RecommendedDifficulty float64           `json:"recommended_difficulty"`
	Error                 any               `json:"error"`  // 暫時用不到
	Total                 int               `json:"total"`  // 搜尋結果總數
	Cursor                any               `json:"cursor"` // 暫時用不到
	CursorString          string            `json:"cursor_string"`
}

type BeatmapSet struct {
	Artist        string `json:"artist"`
	ArtistUnicode string `json:"artist_unicode"`

	// 巢狀物件：封面圖集
	Covers Covers `json:"covers"`

	Creator        string `json:"creator"`
	FavouriteCount int    `json:"favourite_count"`
	GenreID        int    `json:"genre_id"`

	Hype any `json:"hype"` // 暫時用不到

	ID           int    `json:"id"`
	LanguageID   int    `json:"language_id"`
	NSFW         bool   `json:"nsfw"`
	Offset       int    `json:"offset"`
	PlayCount    int    `json:"play_count"`
	PreviewURL   string `json:"preview_url"`
	Source       string `json:"source"`
	Spotlight    bool   `json:"spotlight"`
	Status       string `json:"status"`
	Title        string `json:"title"`
	TitleUnicode string `json:"title_unicode"`

	TrackID *int    `json:"track_id"`
	UserID  int     `json:"user_id"`
	Video   bool    `json:"video"`
	BPM     float64 `json:"bpm"`

	CanBeHyped        bool    `json:"can_be_hyped"`
	DeletedAt         *string `json:"deleted_at"` // 可能是 null
	DiscussionEnabled bool    `json:"discussion_enabled"`
	DiscussionLocked  bool    `json:"discussion_locked"`
	IsScoreable       bool    `json:"is_scoreable"`
	LastUpdated       string  `json:"last_updated"`
	LegacyThreadURL   string  `json:"legacy_thread_url"`

	NominationsSummary any `json:"nominations_summary"` // 暫時用不到

	Ranked        int     `json:"ranked"`
	RankedDate    string  `json:"ranked_date"`
	Rating        float64 `json:"rating"`
	Storyboard    bool    `json:"storyboard"`
	SubmittedDate string  `json:"submitted_date"`
	Tags          string  `json:"tags"`

	Availability Availability `json:"availability"`

	// 這裡就是你原本的 Beatmap 陣列
	Beatmaps []Beatmap `json:"beatmaps"`

	PackTags []string `json:"pack_tags"`
}

type Beatmap struct {
	BeatmapSetID     int     `json:"beatmapset_id"`
	DifficultyRating float64 `json:"difficulty_rating"`
	ID               int     `json:"id"`
	Mode             string  `json:"mode"`
	Status           string  `json:"status"`
	TotalLength      int     `json:"total_length"`
	UserID           int     `json:"user_id"`
	Version          string  `json:"version"` // 難度名稱 (Difficulty Name)

	// --- 難度屬性 (Stats) ---
	Accuracy float64 `json:"accuracy"` // OD (Overall Difficulty)
	AR       float64 `json:"ar"`
	CS       float64 `json:"cs"`
	Drain    float64 `json:"drain"` // HP (Health Drain)
	BPM      float64 `json:"bpm"`

	Convert       bool `json:"convert"`
	CountCircles  int  `json:"count_circles"`
	CountSliders  int  `json:"count_sliders"`
	CountSpinners int  `json:"count_spinners"`

	// --- 時間與狀態 ---
	DeletedAt   *string `json:"deleted_at"` // 可能是 null
	HitLength   int     `json:"hit_length"`
	IsScoreable bool    `json:"is_scoreable"`
	LastUpdated string  `json:"last_updated"`

	// --- 其他資訊 ---
	ModeInt   int    `json:"mode_int"`
	PassCount int    `json:"passcount"`
	PlayCount int    `json:"playcount"`
	Ranked    int    `json:"ranked"`
	URL       string `json:"url"`
	Checksum  string `json:"checksum"`
	MaxCombo  int    `json:"max_combo"`
}

type Covers struct {
	Cover       string `json:"cover"`
	Cover2x     string `json:"cover@2x"`
	Card        string `json:"card"`
	Card2x      string `json:"card@2x"`
	List        string `json:"list"`
	List2x      string `json:"list@2x"`
	Slimcover   string `json:"slimcover"`
	Slimcover2x string `json:"slimcover@2x"`
}

type Availability struct {
	DownloadDisabled bool    `json:"download_disabled"`
	MoreInformation  *string `json:"more_information"` // 可能是 null
}
