package com.github.taoofshawn.hoverDdns.jsonModels;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * Class for importing json data
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class Entries {

	private String id;
	private String content;

	public String getId() {
		return id;
	}

	public void setId(String id) {
		this.id = id;
	}

	public String getContent() {
		return content;
	}

	public void setContent(String content) {
		this.content = content;
	}
}
