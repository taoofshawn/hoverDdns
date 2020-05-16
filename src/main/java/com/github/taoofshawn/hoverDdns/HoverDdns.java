package com.github.taoofshawn.hoverDdns;

import java.net.*;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.HashMap;
import java.util.Map;


import com.github.taoofshawn.hoverDdns.jsonModels.Domains;
import com.github.taoofshawn.hoverDdns.jsonModels.Entries;
import com.github.taoofshawn.hoverDdns.jsonModels.HoverAllDns;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.fasterxml.jackson.databind.ObjectMapper;

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

	}


	private final HttpClient httpClient = HttpClient.newBuilder()
			.version(HttpClient.Version.HTTP_2)
			.cookieHandler(new CookieManager())
			.build();

//	private static HttpRequest.BodyPublisher buildFormDataFromMap(Map<Object, Object> data) {
//		var builder = new StringBuilder();
//		for (Map.Entry<Object, Object> entry : data.entrySet()) {
//			if (builder.length() > 0) {
//				builder.append("&");
//			}
//			builder.append(URLEncoder.encode(entry.getKey().toString(), StandardCharsets.UTF_8));
//			builder.append("=");
//			builder.append(URLEncoder.encode(entry.getValue().toString(), StandardCharsets.UTF_8));
//		}
//		System.out.println(builder.toString());
//		return HttpRequest.BodyPublishers.ofString(builder.toString());
//	}

	/*
	public void post(String uri, String data) throws Exception {
		HttpClient client = HttpClient.newBuilder().build();
		HttpRequest request = HttpRequest.newBuilder()
				.uri(URI.create(uri))
				.POST(BodyPublishers.ofString(data))
				.build();

		HttpResponse<?> response = client.send(request, BodyHandlers.discarding());
		System.out.println(response.statusCode());
		*/


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

		HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());


	}

	/*
	private void sendPost() throws Exception {

		// form parameters
		Map<Object, Object> data = new HashMap<>();
		data.put("username", "abc");
		data.put("password", "123");
		data.put("custom", "secret");
		data.put("ts", System.currentTimeMillis());

		HttpRequest request = HttpRequest.newBuilder()
				.POST(buildFormDataFromMap(data))
				.uri(URI.create("https://httpbin.org/post"))
				.setHeader("User-Agent", "Java 11 HttpClient Bot") // add request header
				.header("Content-Type", "application/x-www-form-urlencoded")
				.header("Content-Type", "application/json")
				.build();

		HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

		// print status code
		System.out.println(response.statusCode());

		// print response body
		System.out.println(response.body());

	} */

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
			exc.printStackTrace();
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

}