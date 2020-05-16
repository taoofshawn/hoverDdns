package com.github.taoofshawn.hoverDdns.jsonModels;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

/**
 * Class for importing json data
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class Domains {

	private List<Entries> entries = null;

	public List<Entries> getEntries() {
		return entries;
	}

	public void setEntries(List<Entries> entries) {
		this.entries = entries;
	}
}
