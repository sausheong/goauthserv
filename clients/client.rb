require 'json'
require 'rest-client'

# Create the following users first

email = "sausheong@gmail.com"
password = "123"

puts "Authenticate user"
json = RestClient.post "http://localhost:3000/authenticate", email: email, password: password
puts "JSON result:"
puts json

puts

puts "Validate session"
RestClient.post("http://localhost:3000/validate", session: JSON.parse(json)["session"]) do |response, request, result| 
  puts response.code
end