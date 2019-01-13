package data

// Post 内容结构
type Post struct {
	ID      int64
	Type    int8
	Title   string
	Content string
	Created int
	Updated int
}

// Posts list
func Posts() ([]*Post, error) {
	rows, err := Conn.Query("SELECT `id`,`type`,`title`,`content`,`updated` FROM posts")
	if err != nil {
		return nil, err
	}
	var res []*Post
	for rows.Next() {
		p := Post{}
		err = rows.Scan(&p.ID, &p.Type, &p.Title, &p.Content, &p.Updated)
		res = append(res, &p)
	}
	return res, nil
}

// Create Post
func (p *Post) Create() (err error) {
	stmt, err := Conn.Prepare("INSERT INTO posts(`type`,`title`,`content`) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer Conn.Close()

	res, err := stmt.Exec(p.Type, p.Title, p.Content)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return
}

// Update Post
func (p *Post) Update() (err error) {
	return
}

// Delete Post
func (p *Post) Delete() (err error) {
	return
}
