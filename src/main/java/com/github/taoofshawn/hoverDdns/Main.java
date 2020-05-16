
package com.github.taoofshawn.hoverDdns;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Main {

	public static void main(String[] args) throws Exception {

		final Logger logger = LoggerFactory.getLogger(Main.class);

		String hoveruser = null;
		String hoverpass = null;
		String hoverid = null;
		String polltime = null;

		if (System.getenv("HOVERUSER") == null) {
			logger.error("missing environment variable: HOVERUSER");
			System.exit(1);
		} else if (System.getenv("HOVERPASS") == null) {
			logger.error("missing environment variable: HOVERPASS");
			System.exit(1);
		} else if (System.getenv("HOVERID") == null) {
			logger.error("missing environment variable: HOVERID");
			System.exit(1);
		} else {
			hoveruser = System.getenv("HOVERUSER");
			hoverpass = System.getenv("HOVERPASS");
			hoverid = System.getenv("HOVERID");
		}

		if (System.getenv("POLLTIME") == null) {
			polltime = "360";
		} else {
			polltime = System.getenv("POLLTIME");
		}

		HoverDdns client = new HoverDdns(hoveruser, hoverpass, hoverid);

//		System.out.println(client.getExternalIP());
		client.getAuth();
//		System.out.println(client.getHoverIP());

		String hoverIP = client.getHoverIP();
		String externalIP = client.getExternalIP();

		if (hoverIP.equals(externalIP)) {
			logger.info(String.format("hover IP does not need to be updated. Hover: %s, Actual: %s", hoverIP, externalIP ));
		} else {
			logger.info(String.format("hover IP needs to be updated. Hover: %s, Actual: %s", hoverIP, externalIP ));
			// update IP here
		}




	}

}