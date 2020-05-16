package com.github.taoofshawn.hoverDdns.jsonModels;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.List;

/**
 * Class for importing json data. This is the top-level class for the full hover dns list
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class HoverAllDns {

	private List<Domains> domains = null;

	public List<Domains> getDomains() {
		return domains;
	}

	public void setDomains(List<Domains> domains) {
		this.domains = domains;
	}
}
