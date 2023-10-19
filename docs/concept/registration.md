# Registration concept

The following text will be shown after an user either clicks on "Login", "Sign-Up" or performs an action that requires login:

> It's time for you to `[login]`.
> 
> Or you can `[create a new account]`.
> 
> Note: If you have many personalities, don't worry, create as many accounts as you like.
> You can even use the same email address, as long as the username is different.
> You can even give one account control over another, so you only need to login with one of them.
> Then you can post comments or create pages as another account.
> Create accounts for your teams, organisations, pets, even robots!

When a user clicks on login, different signup methods are shown. Currently they are username/password, facebook and email link.

> How would you like to login?
> 
> * `[Username and Password]`
> * `[Facebook]`
> * `[Email Link]`
>
> If you have forgotten your password, but you have access to your email, then you can use the `Email Link` method.

When a user clicks on "Sign-Up", the following text is shown instead:

> Please choose a username:
> 
> [______________________]
> 
> * Between 3 and 30 characters.
> * You can use letters, numbers and -_.
> * Just no period (.) at the start.
>
> If this is too restrictive for you, don't worry: Your **display name** can be anything you like.
> 
> `[Next]`

If the username is already taken, the following text appears next to the username field:

> Sorry, that username is already taken.
> You can `[Claim]` it, if the previous registration attempt was incomplete.

Sometimes, a user may see this screen:

> Let's play a quick game to make sure you're not a robot.
> You won't see it often, I promise!
> 
> Please help me spell the word PARKOUR in the correct order.
> 
> (An image with the letters PARKOUR randomly placed and rotated are shown)

If successful, the following text is shown:

> Congratulations, this username is now available.

When a user clicks on "Next", the following text is shown:

> Congratulations, this username is now yours.
> To be able to login, you need to choose at least one login method:
> 
> * `[Password]`
> * `[Facebook]`
> * `[Email Link]`
> 
> `[Continue]`

While clicking on Facebook just uses the Facebook flow, clicking on "Password" shows the following text:

> Please choose a password:
> 
> [______________________]
> 
> Today is not your creative day? What about this one:
> 
> * `urban roofgap aerial reach oak`
> * `east-calm-bush-elbow-vault`
> * `[Try something else]`

If you click on "Email Link", you'll instead see the following screen:

> Please enter your email address and we'll email you a link to login:
> 
> [______________________]
> 
> `[Next]`

Email will be highlighted with (unverified) next to it and the user still needs to choose another
login method or come back later.

When a user finally clicks on "Continue", the following text is shown:

> Tell us more about yourself:
>
> Your name: [______________________]
>
> Very little restrictions here. Just keep it max 100 characters.
> 
> What user type do you represent?
> 
> No matter what you choose here, your experience will always be the same!
> But when other people search or find you, they know who you represent.
> 
> If you want to represent more than one of the following, consider creating more accounts later!
> 
> - **user** *(currently selected)*
>   - Not specified role
> - athlete
>   - An athlete is actively doing Parkour
> - coach
>   - A coach is teaching Parkour to others
> - team
>   - A team has clearly defined members and is working on Parkour projects together
> - group
>   - A group of people who do their training together, but don't want or need to be a team
> - association
>   - An association is organised and offers classes to its members
> - freelancer
>   - A freelancer is a person who is running a business, but is not a company
> - company
>
>   A company is a legal entity that is offering services to its customers
> - school
>
>   A school is an institution that offers education
> - government
>
>   This could be a government agency, a ministry, a city, a state, a country, etc.
> - robot
>
>   A robot is a computer program that is running on a server. While you don't need a user account
>   to interact with most of the API, you can create a user account for your robot, so you can
>   create and modify pages, comments, import and export data programmatically.
> - ~~administrator~~
>
>   ~~An administrator can create and modify more pages than a normal user~~
>
> `Continue`

When a user clicks on "Continue", the user gets redirected to their profile page.
There they will find that they can upload a picture and add a description about themselves.
They'll also see "no comments yet" or "this user has no pages yet", suggesting that they should create some.

An association could create a page about their coaches, how to find them, what they offer, etc.