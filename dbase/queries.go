package dbase

//Query names for function GetQuery
const (
	SUserBySessionID          = "SelectUserBySessionID"
	SUserByEmail              = "SelectUserByEmail"
	SUserByID                 = "SelectUserByID"
	SSessions                 = "SelectSessions"
	STeachers                 = "SelectTeachers"
	STeacherByID              = "SelectTeacherByID"
	STeacherByUserID          = "SelectTeacherByUserID"
	SStudentsByTeacher        = "SelectStudentsByTeacher"
	SelectQuestions           = "SelectQuestions"
	SelectQuestionByID        = "SelectQuestionByID"
	SelectAnswersByQuestionID = "SelectAnswersByQuestionID"
	SLevels                   = "SelectLevels"
	IUser                     = "InsertUser"
	ISession                  = "InsertSession"
	ITeacher                  = "InsertTeacher"
	ILevel                    = "InsertLevel"
	InsertAnswer              = "InsertAnswer"
	ULevel                    = "UpadteLevel"
	USessionsLastActivity     = "UpdateSessionsLastActivity"
	UpdateAnswer              = "UpdateAnswer"
	UpdateAnswerDeletedAt     = "UpdateAnswerDeletedAt"
	UpdateTeacherDeletedAt    = "UpdateTeacherDeletedAt"
	UpdateTeacher             = "UpdateTeacher"
	UpdateUser                = "UpdateUser"
	DSessionByID              = "DeleteSessionByID"
	DSessionByUUID            = "DeleteSessionByUUID"
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
				u.first_name,
				u.last_name,
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
				u.first_name,
				u.last_name,
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
				u.first_name,
				u.last_name,
				u.type,
				u.userpic
			from sessions s
  				left join users u
					on s.user_id = u.id
			where
				s.uuid = $1;`
	case SSessions:
		result = `
			select
				s.id,
				s.uuid,
				s.user_id,
				s.last_activity,
				s.ip,
				s.user_agent
			from sessions s;`
	case STeachers:
		result = `
			select
  				t.id,
				t.user_id,
				t.level_id,
  				t.deleted_at,				
  				l.name,
  				u.email,
  				u.first_name,
  				u.last_name,
  				case when u.userpic is null then 'defaultuserpic.png' else u.userpic end as userpic
			from teachers t
  			left join users u
    			on t.user_id = u.id
  			left join levels l
    			on t.level_id = l.id
			order by
				t.deleted_at desc,
				u.first_name,
				u.last_name;`
	case STeacherByID:
		result = `
			select
  				t.id,
				t.user_id,
				t.level_id,
  				t.deleted_at,				
  				l.name,
  				u.email,
  				u.first_name,
  				u.last_name,
  				case when u.userpic is null then 'defaultuserpic.png' else u.userpic end as userpic
			from teachers t
  			left join users u
    			on t.user_id = u.id
  			left join levels l
    			on t.level_id = l.id
			where
				t.id = $1;`
	case STeacherByUserID:
		result = `
			select
				t.id,
				t.level_id
			from teachers t
			where
				t.user_id = $1;`
	case SStudentsByTeacher:
		result = `
			select
				s.id,
				s.user_id,
				s.level_id
			from students s
			where
				s.teacherid = $1;`
	case SelectQuestions:
		result = `
			select
  				q.id,
  				q.name,
  				q.type,
  				q.score,
  				q.created_at,
  				q.level_id,
  				l.name
			from questions q
  				left join levels l
    				on q.level_id = l.id
			order by
  				l.score,
  				q.type,
				q.created_at;`
	case SelectQuestionByID:
		result = `
			select
  				q.id,
  				q.name,
  				q.type,
  				q.score,
  				q.created_at,
  				q.level_id,
  				l.name
			from questions q
  				left join levels l
    				on q.level_id = l.id
			where
				q.id = $1
			order by
  				l.score,
  				q.type,
				q.created_at;`
	case SelectAnswersByQuestionID:
		result = `
			select
  				a.id,
  				a.name,
 				a.correct,
  				a.created_at,				
  				a.question_id,
				a.deleted_at
			from answers a  				
			where
  				a.question_id = $1
			order by
				a.deleted_at desc,
				a.correct desc,
				a.name`
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
   				first_name,
   				last_name,
   				type,
				userpic)
			values ($1, $2, $3, $4, $5, $6);`
	case ISession:
		result = `
			insert into sessions
				(uuid,
   				user_id,
   				last_activity,
				ip,
				user_agent)
			values ($1, $2, $3, $4, $5);`
	case ITeacher:
		result = `
			insert into teachers
				(user_id,
   				level_id)
			values ($1, $2);`
	case InsertAnswer:
		result = `
			insert into answers
				(name,
   				correct,
				created_at,
				question_id)
			values ($1, $2, $3, $4);`
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
	case UpdateAnswer:
		result = `
			update answers
			set				
				name=$2,
				correct=$3
			where
				id=$1`
	case UpdateAnswerDeletedAt:
		result = `
			update answers
			set		
				correct=false,
				deleted_at=$2				
			where
				id=$1`
	case UpdateTeacherDeletedAt:
		result = `
			update teachers
			set					
				deleted_at=$2				
			where
				id=$1`
	case UpdateTeacher:
		result = `
			update teachers
			set					
				level_id=$2				
			where
				id=$1`
	case UpdateUser:
		result = `
			update users
			set					
				email=$2,
				first_name=$3,				
				last_name=$4,
				userpic=$5				
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
