package main

type SlackRes struct {
	Type      string              `json:"type"`
	User      SlackUser           `json:"user"`
	Token     string              `json:"token"`
	TriggerId string              `json:"trigger_id"`
	View      SlackView           `json:"view"`
	Actions   []SlackSelectOption `json:"actions"`
}

type SlackUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	TeamId   string `json:"team_id"`
}

type SlackView struct {
	Blocks          []SlackBlock `json:"blocks"`
	State           SlackValues  `json:"state"`
	PrivateMetadata string       `json:"private_metadata"`
	CallbackId      string       `json:"callback_id"`
}

type SlackBlock struct {
	BlockId string `json:"block_id"`
}

type SlackValues struct {
	Values map[string]interface{} `json:"values"`
}

type SlackSelectAction struct {
	StaticSelectAction SlackStaticSelect `json:"static_select-action"`
}

type SlackStaticSelect struct {
	SelectedOption SlackSelectOption `json:"selected_option"`
}

type SlackPlainTextInputAction struct {
	PlainTextInputAction SlackSelectOption `json:"plain_text_input-action"`
}

type SlackSelectOption struct {
	Value string `json:"value"`
}
