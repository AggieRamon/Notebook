package main

type ViewResponse struct {
	TriggerId string     `json:"trigger_id"`
	View      SlackModal `json:"view"`
}

type SlackModal struct {
	Type            string       `json:"type"`
	PrivateMetadata string       `json:"private_metadata,omitempty"`
	CallbackId      string       `json:"callback_id,omitempty"`
	Title           SlackText    `json:"title"`
	Submit          SlackText    `json:"submit"`
	Close           SlackText    `json:"close"`
	Blocks          []SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type    string       `json:"type"`
	Element SlackElement `json:"element"`
	Label   SlackText    `json:"label"`
}

type SlackElement struct {
	Type        string        `json:"type"`
	Placeholder *SlackText    `json:"placeholder,omitempty"`
	Options     []SlackOption `json:"options,omitempty"`
	ActionId    string        `json:"action_id,omitempty"`
	Text        string        `json:"text,omitempty"`
	Multiline   bool          `json:"multiline,omitempty"`
}

type SlackText struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

type SlackOption struct {
	Text  SlackText `json:"text"`
	Value string    `json:"value"`
}
