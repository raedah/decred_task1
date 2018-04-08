# decred_task1
a simple website in GoLang that allows for a user to create a new account and to login. Password should be securely stored using bcrypt and salted.

 For a next commit, when the user is logged in you can let them edit their profile. Let them set a Name, Email, Telegram, Skype, WhatsApp, Signal. All fields are optional. Also have them be able to change their password. Another thing, create an invite link, from the users profile page it shows their invite link and they can give it to other users. When someone signs up via invite link the user referral is saved. Adding new used should also save a created date.
 
 You can also allow the user to input their Location. Can just be a text string. Also have them input their timezone, would range from UTC−12:00  to UTC+14   [ref](https://en.wikipedia.org/wiki/List_of_UTC_time_offsets)
 
 no need to put the hour range. Just put the sum hours. Include decription of work done and URL for the commit.
 
 in the login.html file I see it says   email/password.   It should be  username/password
 The username should also be enforced to contain no spaces.
 
 in addition to CreatedDate, lets also add a LastLoginDate and LastModifiedDate for when they update their profile
 
 Allow a user to be set as admin=true. If the user has the admin flag, give them a view to see a list of all users.
 
 Have a messages view, where users can type a msg into a form, similiar to chat. The chat msgs to go the admin. The admin can see a list of unread msgs, and respond to the user msgs.
 
 Have a projects view. Admin can add new project names, and descriptions.  Users can see it but not modify it. Users can select a project to specify they are working on it.
 
 Add worklog view. Users inputs   date, # of hours, description, url reference.
 
 A worklog row would generally be connected to a project.
 
