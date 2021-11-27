package testdata

import (
	"io"
	"log"
	"net/http"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
)

const avatarURL = "https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/daaa1ab5-76fc-477b-9f9a-33d12d9835c6/dds3swl-44008d15-d173-46af-89fa-868e744cdff6.png/v1/fill/w_1280,h_1280,q_80,strp/ferris_argyle_pixelart_by_asoru_dds3swl-fullview.jpg?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7ImhlaWdodCI6Ijw9MTI4MCIsInBhdGgiOiJcL2ZcL2RhYWExYWI1LTc2ZmMtNDc3Yi05ZjlhLTMzZDEyZDk4MzVjNlwvZGRzM3N3bC00NDAwOGQxNS1kMTczLTQ2YWYtODlmYS04NjhlNzQ0Y2RmZjYucG5nIiwid2lkdGgiOiI8PTEyODAifV1dLCJhdWQiOlsidXJuOnNlcnZpY2U6aW1hZ2Uub3BlcmF0aW9ucyJdfQ.O97mCEQAR2F9Az4jx4muSg779uv02znysBBoLN5iRh0"

func MustAvatarPixbuf() *gdkpixbuf.Pixbuf {
	p, err := avatarPixbuf()
	if err != nil {
		log.Panicln("AvatarPixbuf:", err)
	}
	return p
}

func avatarPixbuf() (*gdkpixbuf.Pixbuf, error) {
	r, err := http.Get(avatarURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	pl := gdkpixbuf.NewPixbufLoader()

	_, err = io.Copy(gioutil.PixbufLoaderWriter(pl), r.Body)
	if err != nil {
		return nil, err
	}

	if err := pl.Close(); err != nil {
		return nil, err
	}

	return pl.Pixbuf(), nil
}
