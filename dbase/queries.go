package dbase

func SelectTeachers() string {
	return `select
  				t.id,
				t.user_id,
				t.level_id,
  				t.deleted_at,				
  				l.name,
  				u.email,
  				u.first_name,
  				u.last_name,
  				case when u.userpic is null or u.userpic = '' then 'defaultuserpic.png' else u.userpic end as userpic
			from teachers t
  			left join users u
    			on t.user_id = u.id
  			left join levels l
    			on t.level_id = l.id
			order by
				t.deleted_at desc,
				u.first_name,
				u.last_name;`
}

func SelectTeacherByID() string {
	return `select
  				t.id,
				t.user_id,
				t.level_id,
  				t.deleted_at,				
  				l.name,
  				u.email,
  				u.first_name,
  				u.last_name,
  				case when u.userpic is null or u.userpic = '' then 'defaultuserpic.png' else u.userpic end as userpic
			from teachers t
  			left join users u
    			on t.user_id = u.id
  			left join levels l
    			on t.level_id = l.id
			where
				t.id = $1;`
}

func SelectUserByEmail() string {
	return `select
				u.id,
				u.email,
				u.password,
				u.first_name,
				u.last_name,
				u.type,
				case when u.userpic is null or u.userpic = '' then 'defaultuserpic.png' else u.userpic end as userpic
			from users u
			where
				u.email = $1;`
}

func SelectUserByID() string {
	return `select
				u.id,
				u.email,
				u.password,
				u.first_name,
				u.last_name,
				u.type,
				case when u.userpic is null or u.userpic = '' then 'defaultuserpic.png' else u.userpic end as userpic
			from users u
			where
				u.id = $1;`
}

func SelectUserBySessionID() string {
	return `select 
  				u.id,
				u.email,
				u.password,
				u.first_name,
				u.last_name,
				u.type,
				case when u.userpic is null or u.userpic = '' then 'defaultuserpic.png' else u.userpic end as userpic
			from sessions s
  				left join users u
					on s.user_id = u.id
			where
				s.uuid = $1;`
}

func SelectSessions() string {
	return `select
				s.id,
				s.uuid,
				s.user_id,
				s.last_activity,
				s.ip,
				s.user_agent
			from sessions s;`
}

func SelectTeacherByUserID() string {
	return `select
				t.id,
				t.level_id
			from teachers t
			where
				t.user_id = $1;`
}

func SelectStudentsByTeacher() string {
	return `select
				s.id,
				s.user_id,
				s.level_id
			from students s
			where
				s.teacherid = $1;`
}

func SelectQuestions() string {
	return `select
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
}

func SelectQuestionByID() string {
	return `select
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
}

func SelectAnswersByQuestionID() string {
	return `select
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
				a.name;`
}

func SelectLevels() string {
	return `select
				l.id,
				l.name,
				l.score
			from levels l
			order by
				l.id;`
}

func InsertUser() string {
	return `insert into users
				(email,
   				password,
   				first_name,
   				last_name,
   				type,
				userpic)
			values ($1, $2, $3, $4, $5, $6);`
}

func InsertSession() string {
	return `insert into sessions
				(uuid,
   				user_id,
   				last_activity,
				ip,
				user_agent)
			values ($1, $2, $3, $4, $5);`
}

func InsertTeacher() string {
	return `insert into teachers
				(user_id,
   				level_id)
			values ($1, $2)
			returning id;`
}

func InsertAnswer() string {
	return `insert into answers
				(name,
   				correct,
				created_at,
				question_id)
			values ($1, $2, $3, $4);`
}

func InsertLevel() string {
	return `insert into levels
				(name,
				score)
			values ($1, $2)`
}

func UpdateSessionLastActivityByUuid() string {
	return `update sessions
			set
				last_activity = $2
			where

				uuid=$1;`
}

func UpdateLevel() string {
	return `update levels
			set
				name=$2,
				score=$3
			where
				id=$1;`
}

func UpdateAnswer() string {
	return `update answers
			set				
				name=$2,
				correct=$3
			where
				id=$1;`
}

func UpdateAnswerDeletedAt() string {
	return `update answers
			set		
				correct=false,
				deleted_at=$2				
			where
				id=$1;`
}

func UpdateTeacherDeletedAt() string {
	return `update teachers
			set					
				deleted_at=$2				
			where
				id=$1;`
}

func UpdateTeacher() string {
	return `update teachers
			set					
				level_id=$2				
			where
				id=$1;`
}

func UpdateUser() string {
	return `update users
			set					
				email=$2,
				first_name=$3,				
				last_name=$4,
				userpic=$5				
			where
				id=$1;`
}

func DeleteSessionByID() string {
	return `delete				
			from sessions s
			where
				s.id = $1;`
}

func DeleteSessionByUUID() string {
	return `delete				
			from sessions s
			where
				s.uuid = $1;`
}
