import datetime, json, logging, os, sys, time
import requests

class HoverClient:
	def __init__(self, config):
		self.username = config['HOVERUSER']
		self.password = config['HOVERPASS']
		self.getAuth()
		self.hoverId = config['HOVERID']
		self.hoverIp = self.getCurrentHoverIp()
		self.currentIp = self.getCurrentExternalIp()

	def getAuth(self):
		logging.info('attempting authentication with Hover')
		data = json.dumps({ "username": self.username, "password": self.password })
		loginUrl = 'https://www.hover.com/api/login'
		headers = {'Content-type': 'application/json'}

		r = requests.post(loginUrl, data=data, headers=headers)
		if not r.ok or "hoverauth" not in r.cookies:
			logging.error('could not authenticate with Hover')
			sys.exit(1)
		logging.info('authenticated successfully with Hover')

		self.hoverToken = {"hoverauth": r.cookies["hoverauth"]}
		self.hoverTokenTimestamp = datetime.datetime.now()
		return

	def checkAuth(self):
		lastAuthDelta = datetime.datetime.now() - self.hoverTokenTimestamp
		if lastAuthDelta.total_seconds() > (6 * 60 * 60): #ReAuth after 6 Hours
			logging.info('reauthentication needed')
			self.getAuth()
		return

	def call(self, method, resource, data=None):
		self.checkAuth()
		url = "https://www.hover.com/api/{0}".format(resource)
		logging.info('connecting to Hover')
		r = requests.request(method, url, data=data, cookies=self.hoverToken)
		if not r.ok:
			logging.error('could not connect to Hover!')
		if r.content:
			body = r.json()
			if "succeeded" not in body or body["succeeded"] is not True:
				logging.error('connected to Hover but data is missing!')
			logging.info('request to Hover was successful')
			return body

	def getCurrentHoverIp(self):
		logging.info('checking current IP address setting at Hover')
		allHoverData = self.call("get", "dns")
		for domain in allHoverData.get('domains'):
			for entry in domain['entries']:
				if entry['id'] == self.hoverId:
					return entry['content']
		return 'hoverConnectFailed'

	def getCurrentExternalIp(self):
		logging.info('checking current external IP address')
		return requests.get('https://api.ipify.org').text

	def updateHoverIp(self, newIp):
		logging.info('updating hover IP with {}'.format(newIp))
		return self.call('put', 'dns/' + self.hoverId, {'content': newIp})


if __name__ == "__main__":
	config = {
		'HOVERUSER': os.getenv('HOVERUSER'),
		'HOVERPASS': os.getenv('HOVERPASS'),
		'HOVERID': os.getenv('HOVERID'),
		'POLLTIME': os.getenv('POLLTIME', default='360'),
		'LOGLEVEL': os.getenv('LOGLEVEL', default='INFO')
	}

	logging.basicConfig(format='%(asctime)s - %(message)s', level=config['LOGLEVEL'])

	for key in config:
		if config[key] is None:
			logging.error('missing environment variable: {}'.format(key))
			sys.exit(1)

	client = HoverClient(config)
	while True:
		if client.hoverIp != client.currentIp:
			logging.info(
				'hover IP needs to be updated. Hover: {}, Actual: {}'.format(
					client.hoverIp, client.currentIp))
			try:
				client.updateHoverIp(client.currentIp)
			except:
				logging.error('hover update failed')
		else:
			logging.info(
				'hover IP does not need to be updated. Hover: {}, Actual: {}'.format(
					client.hoverIp, client.currentIp))

		try:
			polltime = int(config['POLLTIME'])
		except:
			polltime = 360
		logging.info('sleeping for {} minutes'.format(polltime))
		time.sleep(polltime*60)  # convert to seconds

		client.hoverIp = client.getCurrentHoverIp()
		client.currentIp = client.getCurrentExternalIp()
