
package com.github.taoofshawn.hoverDdns;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

public class Main {

	public static void main(String[] args) throws Exception {

		final Logger logger = LoggerFactory.getLogger(Main.class);

		String hoveruser = null;
		String hoverpass = null;
		String hoverid = null;
		int polltime;

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


		try {
			polltime = Integer.parseInt(System.getenv("POLLTIME"));

		} catch (Exception exp) {
			polltime = 360;
		}


		HoverDdns client = new HoverDdns(hoveruser, hoverpass, hoverid);

		String hoverIP = "";
		String externalIP = "";

		while (true) {

			try {
				hoverIP = client.getHoverIP();
				externalIP = client.getExternalIP();
			} catch (Exception exc) {
				logger.error("unable to get current IP information");
			}

			if (hoverIP.equals(externalIP)) {
				logger.info(String.format("hover IP does not need to be updated. Hover: %s, Actual: %s", hoverIP, externalIP));
			} else {
				logger.info(String.format("hover IP needs to be updated. Hover: %s, Actual: %s", hoverIP, externalIP));

				try {
					client.updateHoverIP(externalIP);
				} catch (Exception exc) {
					logger.error("hover update failed");
				}
			}

			logger.info(String.format("sleeping for %d minutes", polltime));
			TimeUnit.MINUTES.sleep(polltime);
		}

	}

}