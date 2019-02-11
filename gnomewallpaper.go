package monitor

import (
	"strconv"
	"strings"
	"time"
)

type GnomeWallpaper struct {
	// The name of the directory containing this XML file, if it's not
	// "pixmaps", "images" or "contents". May use the parent of the parent.
	CollectionName string

	// Path is the full path to the XML file
	Path string

	// Config contains the parsed XML. See: gnomexml.go
	Config *GBackground
}

func NewGnomeWallpaper(name string, path string, config *GBackground) *GnomeWallpaper {
	return &GnomeWallpaper{name, path, config}
}

// StartTime returns the timed wallpaper start time, as a time.Time
func (gw *GnomeWallpaper) StartTime() time.Time {
	// gw.Config.StartTime is a struct that contains ints,
	// where the values are directly from the parsed XML.
	st := gw.Config.StartTime
	return time.Date(st.Year, time.Month(st.Month), st.Day, st.Hour, st.Minute, 0, 0, time.Local)
}

func (gw *GnomeWallpaper) Images() []string {
	var filenames []string
	for _, static := range gw.Config.Statics {
		filenames = append(filenames, static.Filename)
	}
	for _, transition := range gw.Config.Transitions {
		filenames = append(filenames, transition.FromFilename)
		filenames = append(filenames, transition.ToFilename)
	}
	return unique(filenames)
}

// String builds a string with various information about this GNOME timed wallpaper
func (gw *GnomeWallpaper) String() string {
	var sb strings.Builder
	sb.WriteString("--- ")
	sb.WriteString(gw.CollectionName)
	sb.WriteString("---\npath\t\t\t= ")
	sb.WriteString(gw.Path)
	sb.WriteString("\nstart time\t\t= ")
	sb.WriteString(gw.StartTime().String())
	sb.WriteString("\nnumber of static tags\t= ")
	sb.WriteString(strconv.Itoa(len(gw.Config.Statics)))
	sb.WriteString("\nnumber of transitions\t= ")
	sb.WriteString(strconv.Itoa(len(gw.Config.Transitions)))
	sb.WriteString("\nuses these images:\n")
	for _, filename := range gw.Images() {
		sb.WriteString("\t" + filename + "\n")
	}
	return sb.String()
}
