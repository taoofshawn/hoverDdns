package com.github.taoofshawn.hoverDdns;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.taoofshawn.hoverDdns.jsonModels.Domains;
import com.github.taoofshawn.hoverDdns.jsonModels.Entries;
import com.github.taoofshawn.hoverDdns.jsonModels.HoverAllDns;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.CookieManager;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.HashMap;
import java.util.Map;

public class HoverDdns {

	private final Logger logger = LoggerFactory.getLogger(HoverDdns.class);

	private final String username;
	private final String password;
	private final String hoverID;
	private HoverAllDns allHoverData;

	public HoverDdns(String envHoverUser, String envHoverPass, String envHoverID) {

		username = envHoverUser;
		password = envHoverPass;
		hoverID = envHoverID;

		try {
			getAuth();
		} catch (Exception exc) {
			logger.error("could not authenticate with Hover");
			System.exit(1);
		}

	}


	private final HttpClient httpClient = HttpClient.newBuilder()
			.version(HttpClient.Version.HTTP_2)
			.cookieHandler(new CookieManager())
			.build();


	public void getAuth() throws Exception {

		logger.info("attempting authentication with Hover");

		Map<Object, Object> auth = new HashMap<>();
		auth.put("username", username);
		auth.put("password", password);

		ObjectMapper objectMapper = new ObjectMapper();
		String requestBody = objectMapper
				.writerWithDefaultPrettyPrinter()
				.writeValueAsString(auth);

		HttpRequest request = HttpRequest.newBuilder()
				.POST(HttpRequest.BodyPublishers.ofString(requestBody))
				.uri(URI.create("https://www.hover.com/api/login"))
				.setHeader("User-Agent", "Java 11 HttpClient Bot")
				.header("Content-Type", "application/json")
				.build();

		httpClient.send(request, HttpResponse.BodyHandlers.ofString());

	}


	public String getExternalIP() throws Exception {
		logger.info("checking current external IP address");

		HttpRequest request = HttpRequest.newBuilder()
				.GET()
				.uri(URI.create("https://api.ipify.org"))
				.setHeader("User-Agent", "Java 11 HttpClient Bot")
				.header("Content-Type", "application/json")
				.build();

		HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

		return (response.body());
	}


	public String getHoverIP() throws Exception {
		logger.info("checking current IP address setting at Hover");

		HttpRequest request = HttpRequest.newBuilder()
				.GET()
				.uri(URI.create("https://www.hover.com/api/dns"))
				.setHeader("User-Agent", "Java 11 HttpClient Bot")
				.build();

		HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

		try {
			ObjectMapper objectMapper = new ObjectMapper();
			allHoverData = objectMapper.readValue(response.body(), HoverAllDns.class);
		} catch (Exception exc) {
			logger.error("unable to parse json response");
		}

		String hoverIP = "hoverIP not found";
		for (Domains domain : allHoverData.getDomains()) {
			for (Entries entry : domain.getEntries()) {
				if (entry.getId().equals(hoverID)) {
					hoverIP = entry.getContent();
				}
			}
		}

		return hoverIP;
	}


	public void updateHoverIP(String hoverIP) throws Exception {
		logger.info(String.format("updating hover IP with %s", hoverIP));

		HttpRequest request = HttpRequest.newBuilder()
				.uri(URI.create(String.format("https://www.hover.com/api/dns/%s", hoverID)))
				.header("Content-Type", "application/x-www-form-urlencoded")
				.PUT(HttpRequest.BodyPublishers.ofString(String.format("content=%s", hoverIP)))
				.build();

		httpClient.send(request, HttpResponse.BodyHandlers.ofString());

	}

}