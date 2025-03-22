export interface User {
    id:                 number
    user_id:            string
    super_user:         boolean
    name:               string
    email:              string
    avatar_url:         string
    joined:             string
    fav_fictions:       Fiction[]
    contributed_fic:    Fiction[]
}

export interface Fiction {
    id:                 number
    contributor_id:     number
    contributor_name:   string
    cover:              string
    title:              string
    subtitle:           string
    author:             string
    artist:             string
    status:             "Ongoing" | "Completed" | "Hiatus" | "Dropped"
    synopsis:           string
    genres:             Genre[]
    chapters:           Chapter[]
    created:            string
}

export interface Genre {
    id:         number
    genre_name: string
}

export interface Chapter {
    fiction_id: number
    id:         number
    title:      string
    content:    string
    created:    string
}

export type FictionForm = Omit<Fiction, "id" | "contributor_id" | "contributor_name" | "created" | "genres" | "chapters">
export type ChapterForm = Omit<Chapter, "fiction_id" | "id" | "created">
