package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/google/uuid"
)
//This is the structure of the User parameters. Id and secret code are created once we use POST request on Postman.
type User struct {
	ID         string
	SecretCode string
	Name       string
	Email      string
	Playlists  []Playlist
}
//This is the structure of the playlist parameters. Id is created once we use POST request on Postman.
type Playlist struct {
	ID     string
	Name   string
	Songs  []Song
	UserID string
}
//This is the structure of the song parameters. Id is created once we use POST request on Postman.
type Song struct {
	ID       string
	Name     string
	Composers string
	MusicURL string
}
//This struct holds the state of the application. It has maps for storing users, playlists, and songs. It also includes a mutex for thread-safe access to the maps.
type MusicListerAPI struct {
	Users     map[string]User
	Playlists map[string]Playlist
	Songs     map[string]Song
	Mutex     sync.RWMutex
}

func NewMusicListerAPI() *MusicListerAPI {
	return &MusicListerAPI{
		Users:     make(map[string]User),
		Playlists: make(map[string]Playlist),
		Songs:     make(map[string]Song),
	}
}
//This function Registers the user and creates the user id and secret id
/*{"ID":"a77aec87-999d-4e34-b30b-5aef6ef7b85a","SecretCode":"ad2c79b6-6289-4508-893c-9fe78b8ea727","Name":"Harsh Mavi","Email":"maviharsh@gmail.com","Playlists":null}
*/
func (api *MusicListerAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	api.Mutex.Lock()
	defer api.Mutex.Unlock()

	// Check if user with the same email already exists
	for _, user := range api.Users {
		if user.Email == newUser.Email {
			http.Error(w, "User with this email already exists", http.StatusBadRequest)
			return
		}
	}

	newUser.ID = generateUniqueID()
	newUser.SecretCode = generateUniqueID()
	api.Users[newUser.SecretCode] = newUser

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
//It is used to login the user using the secret code generated in the register function 
func (api *MusicListerAPI) LoginUser(w http.ResponseWriter, r *http.Request) {
	secretCode := r.URL.Query().Get("secretCode")

	api.Mutex.RLock()
	defer api.Mutex.RUnlock()

	user, exists := api.Users[secretCode]
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
//this function shows the profile with its parameters based on the secretkey. 
func (api *MusicListerAPI) ViewProfile(w http.ResponseWriter, r *http.Request) {
	secretCode := r.URL.Query().Get("secretCode")

	api.Mutex.RLock()
	defer api.Mutex.RUnlock()

	user, exists := api.Users[secretCode]
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
//this function shows all the songs with its parameters based on the paylist id generated while creating the playlist.
func (api *MusicListerAPI) GetAllSongsOfPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.URL.Query().Get("playlistId")

	api.Mutex.RLock()
	defer api.Mutex.RUnlock()

	playlist, exists := api.Playlists[playlistID]
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(playlist.Songs)
	} else {
		http.Error(w, "Playlist not found", http.StatusNotFound)
	}
}
//this function created the playlist using the secret id generated 
func (api *MusicListerAPI) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	secretCode := r.URL.Query().Get("secretCode")

	api.Mutex.Lock()
	defer api.Mutex.Unlock()

	user, exists := api.Users[secretCode]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var newPlaylist Playlist
	err := json.NewDecoder(r.Body).Decode(&newPlaylist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newPlaylist.ID = generateUniqueID()
	newPlaylist.UserID = user.ID
	api.Playlists[newPlaylist.ID] = newPlaylist

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPlaylist)
}
//this function deletes the playlist using the playlist id
func (api *MusicListerAPI) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.URL.Query().Get("playlistId")

	api.Mutex.Lock()
	defer api.Mutex.Unlock()

	_, exists := api.Playlists[playlistID]
	if exists {
		delete(api.Playlists, playlistID)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Playlist not found", http.StatusNotFound)
	}
}
//this function gives the song details using songId
func (api *MusicListerAPI) GetSongDetail(w http.ResponseWriter, r *http.Request) {
	songID := r.URL.Query().Get("songId")

	api.Mutex.RLock()
	defer api.Mutex.RUnlock()

	song, exists := api.Songs[songID]
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(song)
	} else {
		http.Error(w, "Song not found", http.StatusNotFound)
	}
}
//this function is used to add song to the playlist using the playlist id
func (api *MusicListerAPI) addSongToPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.URL.Query().Get("playlistId")

	api.Mutex.Lock()
	defer api.Mutex.Unlock()

	playlist, exists := api.Playlists[playlistID]
	if !exists {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	var newSong Song
	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newSong.ID = generateUniqueID()
	playlist.Songs = append(playlist.Songs, newSong)
	api.Playlists[playlistID] = playlist

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}
//this function created the unique id
func generateUniqueID() string {
	id := uuid.New()
	return id.String()
}

func main() {
	api := NewMusicListerAPI()

	http.HandleFunc("/register", api.RegisterUser)//post request on postman
	http.HandleFunc("/login", api.LoginUser)//get request on postman
	http.HandleFunc("/ViewProfile", api.ViewProfile)//get request on postman
	http.HandleFunc("/getAllSongsOfPlaylist", api.GetAllSongsOfPlaylist)//get request on postman
	http.HandleFunc("/createPlaylist", api.CreatePlaylist)//post request on postman
	http.HandleFunc("/deletePlaylist", api.DeletePlaylist)//delete request on postman
	http.HandleFunc("/getSongDetail", api.GetSongDetail)//get request on postman
	http.HandleFunc("/addSongToPlaylist", api.addSongToPlaylist)//post request on postman

	fmt.Println("Server is running on :8080")// tyhe server is running on the post 8080
	http.ListenAndServe(":8080", nil)
}
