package dbase

const (
	SelectUserBySessionID      = 1001
	SelectUserByEmail          = 1002
	SelectSessions             = 1003
	SelectTeacherByUserID      = 1004
	SelectStudentsByTeacher    = 1005
	SelectLevels               = 1006
	InsertUser                 = 2001
	InsertSession              = 2002
	InsertTeacher              = 2003
	UpdateSessionsLastActivity = 3001
	DeleteSessionByID        = 4001
	DeleteSessionByUUID        = 4002
)

func GetQuery(QryID int) string {

	var result string

	switch QryID {
	case SelectUserByEmail:
		result = `
			select
				u.id,
				u.email,
				u.password,
				u.firstname,
				u.lastname,
				u.type
			from users u
			where
				u.email = $1;`
	case SelectUserBySessionID:
		result = `
			select 
  				u.id,
  				u.firstname,
  				u.lastname,
  				u.type
			from sessions s
  				left join users u
					on s.userid = u.id
			where
				s.uuid = $1;`
	case SelectSessions:
		result = `
			select
				s.id,
				s.uuid,
				s.userid,
				s.lastactivity,
				s.ip,
				s.useragent
			from sessions s;`
	case SelectTeacherByUserID:
		result = `
			select
				t.id,
				t.levelid,
			from teachers t
			where
				t.userid = $1;`
	case SelectStudentsByTeacher:
		result = `
			select
				s.id,
				s.userid,
				s.levelid,
			from students s
			where
				s.teacherid = $1;`
	case SelectLevels:
		result = `
			select
				l.id,
				l.name,
				l.score
			from levels l;`
	case InsertUser:
		result = `
			insert into users
				(email,
   				password,
   				firstname,
   				lastname,
   				type)
			values ($1, $2, $3, $4, $5);`
	case InsertSession:
		result = `
			insert into sessions
				(uuid,
   				userid,
   				lastactivity,
				ip,
				useragent)
			values ($1, $2, $3, $4, $5);`
	case InsertTeacher:
		result = `
			insert into teachers
				(userid,
   				levelid)
			values ($1, $2);`
	case UpdateSessionsLastActivity:
		result = `
			update sessions

			with`
	case DeleteSessionByID:
		result = `
			delete				
			from sessions s
			where
				s.id = $1;`
	case DeleteSessionByUUID:
		result = `
			delete				
			from sessions s
			where
				s.uuid = $1;`
	}

	return result
}
