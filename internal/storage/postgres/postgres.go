package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

type Movie struct {
	Id          int    `json:"movie_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Rating      int    `json:"rating"`
	Actors      []int  `json:"actors"`
}

type Actor struct {
	Id        int    `json:"actor_id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
	Movies    []int  `json:"movies"`
}

const (
	OrderByTitleAsc        = "title_asc"
	OrderByTitleDesc       = "title_desc"
	OrderByReleaseDateAsc  = "release_date_asc"
	OrderByReleaseDateDesc = "release_date_desc"
	OrderByRatingAsc       = "rating_asc"
	OrderByRatingDesc      = "rating_desc"
)

func New(dbUrl string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS actors(
	    actor_id SERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL,
	    gender VARCHAR(10),
		birthdate DATE);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS movies(
	    movie_id SERIAL PRIMARY KEY,
	    title VARCHAR(150) NOT NULL,
	    description VARCHAR(1000) NOT NULL,
	    release_date DATE,
	    rating SMALLINT);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS actor_movie(
	    movie_id INTEGER REFERENCES movies(movie_id),
	    actor_id INTEGER REFERENCES actors(actor_id));
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveActor(name string, gender string, birthdate string) (int, error) {
	const op = "storage.postgres.SaveActor"

	var actorId int
	err := s.db.QueryRow("INSERT INTO actors(name, gender, birthdate) VALUES ($1, $2, $3) RETURNING actor_id",
		name, gender, birthdate).Scan(&actorId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return actorId, nil
}

func (s *Storage) SaveMovie(title string, description string, releaseDate string, rating int, actorsIds []int) (int, error) {
	const op = "storage.postgres.SaveMovie"

	var movieId int
	err := s.db.QueryRow("INSERT INTO movies(title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING movie_id",
		title, description, releaseDate, rating).Scan(&movieId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	err = s.SaveActorMovie(movieId, actorsIds)
	if err != nil {
		if s.DeleteMovie(movieId) != nil {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return movieId, nil
}

func (s *Storage) SaveActorMovie(movieId int, actorsIds []int) error {
	const op = "storage.postgres.SaveActorMovie"

	stmt, err := s.db.Prepare("INSERT INTO actor_movie(movie_id, actor_id) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, actorId := range actorsIds {
		_, err := stmt.Exec(movieId, actorId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateActorName(actorId int, name string) error {
	const op = "storage.postgres.UpdateActorName"

	_, err := s.db.Exec("UPDATE actors SET name=$1 WHERE actor_id=$2", name, actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateActorGender(actorId int, gender string) error {
	const op = "storage.postgres.UpdateActorGender"

	_, err := s.db.Exec("UPDATE actors SET gender=$1 WHERE actor_id=$2", gender, actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateActorBirthdate(actorId int, birthdate string) error {
	const op = "storage.postgres.UpdateActorBirthdate"

	_, err := s.db.Exec("UPDATE actors SET birthdate=$1 WHERE actor_id=$2", birthdate, actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateMovieTitle(movieId int, title string) error {
	const op = "storage.postgres.UpdateMovieTitle"

	_, err := s.db.Exec("UPDATE movies SET title=$1 WHERE movie_id=$2", title, movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateMovieDescription(movieId int, description string) error {
	const op = "storage.postgres.UpdateMovieDescription"

	_, err := s.db.Exec("UPDATE movies SET description=$1 WHERE movie_id=$2", description, movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateMovieReleaseDate(movieId int, releaseDate string) error {
	const op = "storage.postgres.UpdateMovieReleaseDate"

	_, err := s.db.Exec("UPDATE movies SET release_date=$1 WHERE movie_id=$2", releaseDate, movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateMovieRating(movieId int, rating int) error {
	const op = "storage.postgres.UpdateMovieRating"

	_, err := s.db.Exec("UPDATE movies SET rating=$1 WHERE movie_id=$2", rating, movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteActor(actorId int) error {
	const op = "storage.postgres.DeleteActor"

	_, err := s.db.Exec("DELETE FROM actor_movie WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = s.db.Exec("DELETE FROM actors WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovie(movieId int) error {
	const op = "storage.postgres.DeleteMovie"

	_, err := s.db.Exec("DELETE FROM actor_movie WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM movies WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteActorMovie(movieId int, actorsIds []int) error {
	const op = "storage.postgres.DeleteActorMovie"

	stmt, err := s.db.Prepare("DELETE FROM actor_movie WHERE movie_id=$1 AND actor_id=$2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, actorId := range actorsIds {
		_, err := stmt.Exec(movieId, actorId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetActor(actorId int) (Actor, error) {
	const op = "storage.postgres.GetActor"

	var actor Actor
	err := s.db.QueryRow(`SELECT actor_id, name, gender, to_char(birthdate, 'YYYY-MM-DD') 
								FROM actors WHERE actor_id=$1`, actorId).
		Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthdate)
	if err != nil {
		return Actor{}, fmt.Errorf("%s: %w", op, err)
	}

	movies, err := s.GetMoviesByActor(actorId)
	if err != nil {
		return Actor{}, fmt.Errorf("%s: %w", op, err)
	}
	actor.Movies = movies

	return actor, nil
}

func (s *Storage) GetMovie(movieId int) (Movie, error) {
	const op = "storage.postgres.GetMovie"

	var movie Movie
	err := s.db.QueryRow(`SELECT movie_id, title, description, to_char(release_date, 'YYYY-MM-DD'), rating
								FROM movies WHERE movie_id=$1`, movieId).
		Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
	if err != nil {
		return Movie{}, fmt.Errorf("%s: %w", op, err)
	}

	actors, err := s.GetActorsByMovie(movieId)
	if err != nil {
		return Movie{}, fmt.Errorf("%s: %w", op, err)
	}
	movie.Actors = actors

	return movie, nil
}

func (s *Storage) GetActorsByMovie(movieId int) ([]int, error) {
	const op = "storage.postgres.GetActorsByMovie"

	rows, err := s.db.Query("SELECT actor_id FROM actor_movie WHERE movie_id=$1", movieId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var actors []int
	for rows.Next() {
		var actorId int
		err = rows.Scan(&actorId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		actors = append(actors, actorId)
	}

	return actors, nil
}

func (s *Storage) GetMoviesByActor(actorId int) ([]int, error) {
	const op = "storage.postgres.GetMoviesByActor"

	rows, err := s.db.Query("SELECT movie_id FROM actor_movie WHERE actor_id=$1", actorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var movies []int
	for rows.Next() {
		var movieId int
		err = rows.Scan(&movieId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		movies = append(movies, movieId)
	}

	return movies, nil
}

func (s *Storage) GetMovies(sortBy string) ([]Movie, error) {
	const op = "storage.postgres.GetMovies"

	var movies []Movie
	var orderBy string

	switch sortBy {
	case OrderByTitleAsc:
		orderBy = "title ASC"
	case OrderByTitleDesc:
		orderBy = "title DESC"
	case OrderByReleaseDateAsc:
		orderBy = "release_date ASC"
	case OrderByReleaseDateDesc:
		orderBy = "release_date DESC"
	case OrderByRatingAsc:
		orderBy = "rating ASC"
	case OrderByRatingDesc:
		orderBy = "rating DESC"
	default:
		orderBy = "rating DESC"
	}

	query := fmt.Sprintf(`SELECT movie_id, title, description, to_char(release_date, 'YYYY-MM-DD') as release_date_f, rating 
								 FROM movies
								 ORDER BY %s`, orderBy)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actors, err := s.GetActorsByMovie(movie.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movie.Actors = actors

		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *Storage) GetActors() ([]Actor, error) {
	const op = "storage.postgres.GetActors"

	var actors []Actor

	rows, err := s.db.Query("SELECT actor_id, name, gender, to_char(birthdate, 'YYYY-MM-DD') FROM actors")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var actor Actor
		err = rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthdate)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies, err := s.GetMoviesByActor(actor.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actor.Movies = movies
		actors = append(actors, actor)
	}

	return actors, nil
}

func (s *Storage) GetMoviesBySearchRequest(searchRequest string) ([]Movie, error) {
	const op = "storage.postgres.GetMovieBySearchRequest"

	var movies []Movie

	rows, err := s.db.Query(`SELECT DISTINCT m.movie_id, m.title, m.description, to_char(m.release_date, 'YYYY-MM-DD'), m.rating 
								   FROM movies m 
    							   JOIN actor_movie am ON m.movie_id = am.movie_id 
								   JOIN actors a ON am.actor_id = a.actor_id
								   WHERE m.title LIKE $1 or a.name LIKE $1`, "%"+searchRequest+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actors, err := s.GetActorsByMovie(movie.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movie.Actors = actors

		movies = append(movies, movie)
	}

	return movies, nil
}
