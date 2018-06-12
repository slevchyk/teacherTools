package dbase

const (
	S_UserBySessionID      = "SelectUserBySessionID"
	S_UserByEmail          = "SelectUserByEmail"
	S_UserByID             = "SelectUserByID"
	S_Sessions             = "SelectSessions"
	S_Teachers				= "SelectTeachers"
	S_TeacherByID			= "SelectTeacherByID"
	S_TeacherByUserID      = "SelectTeacherByUserID"
	S_StudentsByTeacher    = "SelectStudentsByTeacher"
	S_Levels               = "SelectLevels"
	I_User                 = "InsertUser"
	I_Session              = "InsertSession"
	I_Teacher              = "InsertTeacher"
	I_Level				= "InsertLevel"
	U_Level = "UpadteLevel"
	U_SessionsLastActivity = "UpdateSessionsLastActivity"
	D_SessionByID          = "DeleteSessionByID"
	D_SessionByUUID        = "DeleteSessionByUUID"
)

func GetQuery(QryID string) string {

	var result string

	switch QryID {
	case S_UserByEmail:
		result = `
			select
				u.id,
				u.email,
				u.password,
				u.firstname,
				u.lastname,
				u.type,
				u.userpic
			from users u
			where
				u.email = $1;`
	case S_UserByID:
		result = `
			select
				u.id,
				u.email,
				u.password,
				u.firstname,
				u.lastname,
				u.type,
				u.userpic
			from users u
			where
				u.id = $1;`
	case S_UserBySessionID:
		result = `
			select 
  				u.id,
				u.email,
				u.password,
				u.firstname,
				u.lastname,
				u.type,
				u.userpic
			from sessions s
  				left join users u
					on s.userid = u.id
			where
				s.uuid = $1;`
	case S_Sessions:
		result = `
			select
				s.id,
				s.uuid,
				s.userid,
				s.lastactivity,
				s.ip,
				s.useragent
			from sessions s;`
	case S_Teachers:
		result = `
			select
  				t.id,
  				t.active,
				t.levelid,
  				l.name,
  				u.email,
  				u.firstname,
  				u.lastname,
  				case when u.userpic is null then 'defaultuserpic.png' end as userpic
			from teachers t
  			left join users u
    			on t.userid = u.id
  			left join levels l
    			on t.levelid = l.id;`
	case S_TeacherByID:
		result = `
			select
  				t.id,
  				t.active,
				t.levelid,
  				l.name,
  				u.email,
  				u.firstname,
  				u.lastname,
  				case when u.userpic is null then 'defaultuserpic.png' end as userpic
			from teachers t
  			left join users u
    			on t.userid = u.id
  			left join levels l
    			on t.levelid = l.id
			where
				t.id = $1;`
	case S_TeacherByUserID:
		result = `
			select
				t.id,
				t.levelid,
			from teachers t
			where
				t.userid = $1;`
	case S_StudentsByTeacher:
		result = `
			select
				s.id,
				s.userid,
				s.levelid,
			from students s
			where
				s.teacherid = $1;`
	case S_Levels:
		result = `
			select
				l.id,
				l.name,
				l.score
			from levels l
			order by
				l.id;`
	case I_User:
		result = `
			insert into users
				(email,
   				password,
   				firstname,
   				lastname,
   				type,
				userpic)
			values ($1, $2, $3, $4, $5, $6);`
	case I_Session:
		result = `
			insert into sessions
				(uuid,
   				userid,
   				lastactivity,
				ip,
				useragent)
			values ($1, $2, $3, $4, $5);`
	case I_Teacher:
		result = `
			insert into teachers
				(userid,
   				levelid)
			values ($1, $2);`
	case I_Level:
		result = `
			insert into levels
				(name,
				score)
			values ($1, $2)`
	case U_SessionsLastActivity:
		result = `
			update sessions

			with`
	case U_Level:
		result = `
			update levels
			set
				name=$2,
				score=$3
			where
				id=$1`
	case D_SessionByID:
		result = `
			delete				
			from sessions s
			where
				s.id = $1;`
	case D_SessionByUUID:
		result = `
			delete				
			from sessions s
			where
				s.uuid = $1;`
	}

	return result
}
