import unittest
import subprocess
import requests # Note, you may need to install this package via pip (or pip3)

PORT = 8080
localhost = "localhost" # windows toolbox users will again want to make this the docker machine's ip adress

class Client():

	def putKey(self, key, value, port):
		result = requests.put('http://%s:%s/kv-store/keys/%s'%(localhost, str(port), key), 
							json={'value':value},
							headers = {"Content-Type": "application/json"})
		return self.formatResult(result)
	
	def getKey(self, key, port):
		result = requests.get('http://%s:%s/kv-store/keys/%s'%(localhost, str(port), key),
							headers = {"Content-Type": "application/json"})
		return self.formatResult(result)
	
	def deleteKey(self, key, port):
		result = requests.delete('http://%s:%s/kv-store/keys/%s'%(localhost, str(port), key),
							headers = {"Content-Type": "application/json"})
		return self.formatResult(result)

	def getKeyCount(self, port):
		result = requests.get('http://%s:%s/kv-store/key-count'%(localhost, str(port)),
							headers = {"Content-Type": "application/json"})
		return self.formatResult(result)

	def putViewChange(self, views, port):
		result = requests.put('http://%s:%s/kv-store/view-change'%(localhost, str(port)),
							json={"view": view},
							headers = {"Content-Type": "application/json"})
		return self.formatResult(result)
	# this just turns the requests result object into a simplified json object 
	# containing only fields I care about 
	def formatResult(self, result):
		status_code = result.status_code
		result = result.json()

		if result != None:			
			jsonKeys = ["message", "replaced", "error", "doesExist", "value"]
			result = {k:result[k] for k in jsonKeys if k in result}

			result["status_code"] = status_code
		else:
			result = {"status_code": status_code}

		return result

client = Client()

#### Expected Responses:
put_success = { 	"message":		"Added successfully",
						"replaced": 	False,
						"address":"changeme",
						"status_code":	201}

put_success_no_addr = { 	"message":		"Added successfully",
						"replaced": 	False,
						"status_code":	201}

put_error_no_key = {	"error":	"Value is missing",
							"message":	"Error in PUT",
						"status_code":	400}

put_error_longKey = {"error":	"Key is too long",
							"message":	"Error in PUT",
						"status_code":	400}

update_success = {"message":		"Updated successfully",
						"replaced":		True,
						"status_code":	200}

update_fail_no_key = put_error_no_key

get_success = {	"doesExist":	True,
						"message":		"Retrieved successfully",
						"value":		"changme",
						"status_code":	200}

get_no_key = {	"doesExist":	False,
						"error":		"Key does not exist",
						"message":		"Error in GET",
						"status_code":	404}

del_success = {	"doesExist":	True,
						"message":		"Deleted successfully",
						"address":"changeme",
						"status_code":	200}

del_success_no_addr = {	"doesExist":	True,
						"message":		"Deleted successfully",
						"status_code":	200}

del_no_key = {	"doesExist":	False,
						"error":		"Key does not exist",
						"message":		"Error in DELETE",
						"status_code":	404}

get_key_count = {
	"message":"Key count retrieved successfully","key-count": "changme"
}

put_view_change = {
    "message": "View change successful",
    "shards" : "changme - array of messages",
}


class TestHW1(unittest.TestCase):

### Add New Keys
	# add a new key
	def test_add_1(self):
		result = client.putKey("Test", "a friendly string", PORT)

		#expected = put_success.copy()
		#expected["address"] = result["address"]
		expected = put_success_no_addr

		self.assertEqual(result, expected)

## Get Key Values
	# add and get
	def test_get_2(self):
		key = "AKey"
		value = "a different friendly string"

		result = client.putKey(key, value, PORT)

		self.assertEqual(result["status_code"], put_success["status_code"], 
			msg="add key: failed add, cannot continue test\n%s\n"%result)


		result = client.getKey(key, PORT)
		expected = get_success.copy()
		expected["value"] = value

		self.assertEqual(result, expected)

# ### Update Keys
	# add then update
	def test_update_1(self):
		key = "AValueToUpdate!"

		result = client.putKey(key, "one, one, one, one!", PORT)

		self.assertEqual(result["status_code"], put_success["status_code"], 
			msg="add key: failed add, cannot continue test\n%s\n"%result)
		

		result = client.putKey(key, "two, three, four!", PORT)

		self.assertEqual(result, update_success)

### Delete Keys
	# add and delete
	def test_del_1(self):
		key = "keyToDelete"

		result = client.putKey(key, "delete, delete, delete!", PORT)

		self.assertEqual(result["status_code"], put_success["status_code"], 
			msg="add key: failed add, cannot continue test\n%s\n"%result)

		result = client.deleteKey(key, PORT)

		#expected = del_success.copy()
		#expected["address"] = result["address"]

		self.assertEqual(result, del_success_no_addr)

### Key Count
	# add and delete
	def test_key_count_1(self):
		key = "keyToDelete"

		result = client.get_key_count(PORT)

		self.assertEqual(result["status_code"], put_success["status_code"], 
			msg="add key: failed add, cannot continue test\n%s\n"%result)

		result = client.deleteKey(key, PORT)

		#expected = del_success.copy()
		#expected["address"] = result["address"]

		self.assertEqual(result, del_success_no_addr)


if __name__ == '__main__':
	unittest.main()
