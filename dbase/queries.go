package dbase

//Query names for function GetQuery
const (
	SUserBySessionID      = "SelectUserBySessionID"
	SUserByEmail          = "SelectUserByEmail"
	SUserByID             = "SelectUserByID"
	SSessions             = "SelectSessions"
	STeachers             = "SelectTeachers"
	STeacherByID          = "SelectTeacherByID"
	STeacherByUserID      = "SelectTeacherByUserID"
	SStudentsByTeacher    = "SelectStudentsByTeacher"
	SelectQuestions = "SelectQuestions"
	SLevels               = "SelectLevels"
	IUser                 = "InsertUser"
	ISession              = "InsertSession"
	ITeacher              = "InsertTeacher"
	ILevel                = "InsertLevel"
	ULevel                = "UpadteLevel"
	USessionsLastActivity = "UpdateSessionsLastActivity"
	DSessionByID          = "DeleteSessionByID"
	DSessionByUUID        = "DeleteSessionByUUID"
)

//GetQuery function return query text by query name (u can use const from this pkg)
func GetQuery(QryID string) string {

	var result string

	switch QryID {
	case SUserByEmail:
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
	case SUserByID:
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
	case SUserBySessionID:
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
	case SSessions:
		result = `
			select
				s.id,
				s.uuid,
				s.userid,
				s.lastactivity,
				s.ip,
				s.useragent
			from sessions s;`
	case STeachers:
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
	case STeacherByID:
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
	case STeacherByUserID:
		result = `
			select
				t.id,
				t.levelid,
			from teachers t
			where
				t.userid = $1;`
	case SStudentsByTeacher:
		result = `
			select
				s.id,
				s.userid,
				s.levelid,
			from students s
			where
				s.teacherid = $1;`
	case SelectQuestions:
		result = `
			select
  				q.id,
  				q.question,
  				q.type,
  				q.score,
  				q.datecreated,
  				q.levelid,
  				l.name
			from questions q
  				left join levels l
    				on q.levelid = l.id
			order by
  				l.score,
  				q.type,
				q.datecreated;`
	case SLevels:
		result = `
			select
				l.id,
				l.name,
				l.score
			from levels l
			order by
				l.id;`
	case IUser:
		result = `
			insert into users
				(email,
   				password,
   				firstname,
   				lastname,
   				type,
				userpic)
			values ($1, $2, $3, $4, $5, $6);`
	case ISession:
		result = `
			insert into sessions
				(uuid,
   				userid,
   				lastactivity,
				ip,
				useragent)
			values ($1, $2, $3, $4, $5);`
	case ITeacher:
		result = `
			insert into teachers
				(userid,
   				levelid)
			values ($1, $2);`
	case ILevel:
		result = `
			insert into levels
				(name,
				score)
			values ($1, $2)`
	case USessionsLastActivity:
		result = `
			update sessions

			with`
	case ULevel:
		result = `
			update levels
			set
				name=$2,
				score=$3
			where
				id=$1`
	case DSessionByID:
		result = `
			delete				
			from sessions s
			where
				s.id = $1;`
	case DSessionByUUID:
		result = `
			delete				
			from sessions s
			where
				s.uuid = $1;`
	}

	return result
}
