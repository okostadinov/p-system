package main

import "net/http"

type Flash struct {
	Content string
	Type    flashType
}

type flashType int

const (
	FlashTypeSuccess flashType = iota
	FlashTypeWarning
	FlashTypeDanger
)

// associated the flash type const with its string equivalent
func (ft flashType) String() string {
	return [...]string{"success", "warning", "danger"}[ft]
}

// adds a flash message to the current session
func (app *application) setFlash(w http.ResponseWriter, r *http.Request, content string, flashType flashType) error {
	session, err := app.store.Get(r, "session")
	if err != nil {
		return err
	}

	flash := Flash{Content: content, Type: flashType}
	session.AddFlash(flash)
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

// pops the flash message from the current session and returns it
func (app *application) popFlash(w http.ResponseWriter, r *http.Request) Flash {
	flash := &Flash{}

	session, _ := app.store.Get(r, "session")

	if flashes := session.Flashes(); len(flashes) > 0 {
		flash = flashes[0].(*Flash)
		session.Save(r, w)
	}

	return *flash
}
