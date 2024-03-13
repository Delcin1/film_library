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
	Id          int
	Title       string
	Description string
	ReleaseDate string
	Rating      int
	Actors      []Actor
}

type Actor struct {
	Id        int
	Name      string
	Gender    string
	Birthdate string
	Movies    []Movie
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

func (s *Storage) SaveActor(name string, gender string, birthdate string) error {
	const op = "storage.postgres.SaveActor"

	_, err := s.db.Exec("INSERT INTO actors(name, gender, birthdate) VALUES ($1, $2, $3)", name, gender, birthdate)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveMovie(title string, description string, releaseDate string, rating int, actorsIds []int) error {
	const op = "storage.postgres.SaveMovie"

	rows, err := s.db.Exec("INSERT INTO movies(title, description, release_date, rating) VALUES ($1, $2, $3, $4)", title, description, releaseDate, rating)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	movieId, err := rows.LastInsertId()
	return s.SaveActorMovie(movieId, actorsIds)
}

func (s *Storage) SaveActorMovie(movieId int64, actorsIds []int) error {
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

	_, err := s.db.Exec("DELETE FROM actors WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM actor_movie WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovie(movieId int) error {
	const op = "storage.postgres.DeleteMovie"

	_, err := s.db.Exec("DELETE FROM movies WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM actor_movie WHERE movie_id=$1", movieId)
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

	return nil
}

func (s *Storage) DeleteActorName(actorId int) error {
	const op = "storage.postgres.DeleteActorName"

	_, err := s.db.Exec("UPDATE actors SET name=NULL WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteActorGender(actorId int) error {
	const op = "storage.postgres.DeleteActorGender"

	_, err := s.db.Exec("UPDATE actors SET gender=NULL WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteActorBirthdate(actorId int) error {
	const op = "storage.postgres.DeleteActorBirthdate"

	_, err := s.db.Exec("UPDATE actors SET birthdate=NULL WHERE actor_id=$1", actorId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovieTitle(movieId int) error {
	const op = "storage.postgres.DeleteMovieTitle"

	_, err := s.db.Exec("UPDATE movies SET title=NULL WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovieDescription(movieId int) error {
	const op = "storage.postgres.DeleteMovieDescription"

	_, err := s.db.Exec("UPDATE movies SET description=NULL WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovieReleaseDate(movieId int) error {
	const op = "storage.postgres.DeleteMovieReleaseDate"

	_, err := s.db.Exec("UPDATE movies SET release_date=NULL WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovieRating(movieId int) error {
	const op = "storage.postgres.DeleteMovieRating"

	_, err := s.db.Exec("UPDATE movies SET rating=NULL WHERE movie_id=$1", movieId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetActorsByMovie(movieId int) ([]Actor, error) {
	const op = "storage.postgres.GetActorsByMovie"

	rows, err := s.db.Query("SELECT * FROM actors WHERE actor_id IN (SELECT actor_id FROM actor_movie WHERE movie_id=$1)", movieId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var actors []Actor
	for rows.Next() {
		var actor Actor
		err = rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.Birthdate)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		actors = append(actors, actor)
	}

	return actors, nil
}

func (s *Storage) GetMoviesByActor(actorId int) ([]Movie, error) {
	const op = "storage.postgres.GetMoviesByActor"

	rows, err := s.db.Query("SELECT * FROM movies WHERE movie_id IN (SELECT movie_id FROM actor_movie WHERE actor_id=$1)", actorId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		movies = append(movies, movie)
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

	rows, err := s.db.Query("SELECT * FROM movies ORDER BY $1", orderBy)
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

func (s *Storage) GetActors(sortBy string) ([]Actor, error) {
	const op = "storage.postgres.GetActors"

	var actors []Actor

	rows, err := s.db.Query("SELECT * FROM actors")
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

func (s *Storage) GetMovieBySearchRequest(searchRequest string) ([]Movie, error) {
	const op = "storage.postgres.GetMovieBySearchRequest"

	var movies []Movie

	rows, err := s.db.Query(`SELECT DISTINCT m.* 
								   FROM movies m 
    							   JOIN actor_movie am ON m.movie_id = am.movie_id 
								   JOIN actors a ON am.actor_id = a.actor_id
								   WHERE m.title LIKE $1 or a.name LIKE $1`, searchRequest)
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
