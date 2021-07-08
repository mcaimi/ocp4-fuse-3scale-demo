package com.redhat.db;

import java.time.LocalDateTime;

public class PostEntity {

    private Integer userId;
	private String title;
    private String body;
    private LocalDateTime created;

    public PostEntity() {
	}

    public PostEntity(Integer userId, String title, String body, LocalDateTime created) {
		this.userId = userId;
		this.title = title;
		this.body = body;
		this.created = created;
	}

	@Override
    public String toString() {
        return "PostEntity{" +
                "userId:" + userId +
                ", title:'" + title + '\'' +
                ", body:'" + body + '\'' +
                ", created:" + created +
                '}';
    }

    public Integer getUserId() {
        return userId;
    }

    public void setUserId(Integer userId) {
        this.userId = userId;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getBody() {
        return body;
    }

    public void setBody(String body) {
        this.body = body;
    }

    public LocalDateTime getCreated() {
		return created;
	}

	public void setCreated(LocalDateTime created) {
		this.created = created;
	}
}
