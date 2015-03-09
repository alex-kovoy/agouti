package api

import (
	"encoding/base64"
	"errors"

	"github.com/sclevine/agouti/api/internal/bus"
)

type Session struct {
	Bus
}

type Bus interface {
	Send(method, endpoint string, body, result interface{}) error
}

func Open(url string, capabilities map[string]interface{}) (*Session, error) {
	busClient, err := bus.Connect(url, capabilities)
	if err != nil {
		return nil, err
	}
	return &Session{busClient}, nil
}

func (s *Session) Delete() error {
	return s.Send("DELETE", "", nil, nil)
}

func (s *Session) GetElement(selector Selector) (*Element, error) {
	var result struct{ Element string }

	if err := s.Send("POST", "element", selector, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetElements(selector Selector) ([]*Element, error) {
	var results []struct{ Element string }

	if err := s.Send("POST", "elements", selector, &results); err != nil {
		return nil, err
	}

	elements := []*Element{}
	for _, result := range results {
		elements = append(elements, &Element{result.Element, s})
	}

	return elements, nil
}

func (s *Session) GetActiveElement() (*Element, error) {
	var result struct{ Element string }

	if err := s.Send("POST", "element/active", nil, &result); err != nil {
		return nil, err
	}

	return &Element{result.Element, s}, nil
}

func (s *Session) GetWindow() (*Window, error) {
	var windowID string
	if err := s.Send("GET", "window_handle", nil, &windowID); err != nil {
		return nil, err
	}
	return &Window{windowID, s}, nil
}

func (s *Session) GetWindows() ([]*Window, error) {
	var windowsID []string
	if err := s.Send("GET", "window_handles", nil, &windowsID); err != nil {
		return nil, err
	}

	var windows []*Window
	for _, windowID := range windowsID {
		windows = append(windows, &Window{windowID, s})
	}
	return windows, nil
}

func (s *Session) SetWindow(window *Window) error {
	request := struct {
		Name string `json:"name"`
	}{window.ID}

	return s.Send("POST", "window", request, nil)
}

func (s *Session) SetWindowByName(name string) error {
	request := struct {
		Name string `json:"name"`
	}{name}

	return s.Send("POST", "window", request, nil)
}

func (s *Session) DeleteWindow() error {
	if err := s.Send("DELETE", "window", nil, nil); err != nil {
		return err
	}
	return nil
}

func (s *Session) GetCookies() ([]*Cookie, error) {
	var cookies []*Cookie
	if err := s.Send("GET", "cookie", nil, &cookies); err != nil {
		return nil, err
	}
	return cookies, nil
}

func (s *Session) SetCookie(cookie *Cookie) error {
	if cookie == nil {
		return errors.New("nil cookie is invalid")
	}
	request := struct {
		Cookie *Cookie `json:"cookie"`
	}{cookie}

	return s.Send("POST", "cookie", request, nil)
}

func (s *Session) DeleteCookie(cookieName string) error {
	return s.Send("DELETE", "cookie/"+cookieName, nil, nil)
}

func (s *Session) DeleteCookies() error {
	return s.Send("DELETE", "cookie", nil, nil)
}

func (s *Session) GetScreenshot() ([]byte, error) {
	var base64Image string

	if err := s.Send("GET", "screenshot", nil, &base64Image); err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(base64Image)
}

func (s *Session) GetURL() (string, error) {
	var url string
	if err := s.Send("GET", "url", nil, &url); err != nil {
		return "", err
	}

	return url, nil
}

func (s *Session) SetURL(url string) error {
	request := struct {
		URL string `json:"url"`
	}{url}

	return s.Send("POST", "url", request, nil)
}

func (s *Session) GetTitle() (string, error) {
	var title string
	if err := s.Send("GET", "title", nil, &title); err != nil {
		return "", err
	}

	return title, nil
}

func (s *Session) GetSource() (string, error) {
	var source string
	if err := s.Send("GET", "source", nil, &source); err != nil {
		return "", err
	}

	return source, nil
}

func (s *Session) DoubleClick() error {
	return s.Send("POST", "doubleclick", nil, nil)
}

func (s *Session) MoveTo(region *Element, offset Offset) error {
	request := map[string]interface{}{}

	if region != nil {
		// TODO: return error if not an element
		request["element"] = region.ID
	}

	if offset != nil {
		if xoffset, present := offset.x(); present {
			request["xoffset"] = xoffset
		}

		if yoffset, present := offset.y(); present {
			request["yoffset"] = yoffset
		}
	}

	return s.Send("POST", "moveto", request, nil)
}

func (s *Session) Frame(frame *Element) error {
	var elementID interface{}

	if frame != nil {
		elementID = struct {
			Element string `json:"ELEMENT"`
		}{frame.ID}
	}

	request := struct {
		ID interface{} `json:"id"`
	}{elementID}

	return s.Send("POST", "frame", request, nil)
}

func (s *Session) FrameParent() error {
	return s.Send("POST", "frame/parent", nil, nil)
}

func (s *Session) Execute(body string, arguments []interface{}, result interface{}) error {
	if arguments == nil {
		arguments = []interface{}{}
	}

	request := struct {
		Script string        `json:"script"`
		Args   []interface{} `json:"args"`
	}{body, arguments}

	if err := s.Send("POST", "execute", request, result); err != nil {
		return err
	}

	return nil
}

func (s *Session) Forward() error {
	return s.Send("POST", "forward", nil, nil)
}

func (s *Session) Back() error {
	return s.Send("POST", "back", nil, nil)
}

func (s *Session) Refresh() error {
	return s.Send("POST", "refresh", nil, nil)
}

func (s *Session) GetAlertText() (string, error) {
	var text string
	if err := s.Send("GET", "alert_text", nil, &text); err != nil {
		return "", err
	}
	return text, nil
}

func (s *Session) SetAlertText(text string) error {
	request := struct {
		Text string `json:"text"`
	}{text}
	return s.Send("POST", "alert_text", request, nil)
}

func (s *Session) AcceptAlert() error {
	return s.Send("POST", "accept_alert", nil, nil)
}

func (s *Session) DismissAlert() error {
	return s.Send("POST", "dismiss_alert", nil, nil)
}

func (s *Session) NewLogs(logType string) ([]Log, error) {
	request := struct {
		Type string `json:"type"`
	}{logType}

	var logs []Log
	if err := s.Send("POST", "log", request, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Session) GetLogTypes() ([]string, error) {
	var types []string
	if err := s.Send("GET", "log/types", nil, &types); err != nil {
		return nil, err
	}
	return types, nil
}
