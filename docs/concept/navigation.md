Actions that can be performed on this website are mostly viewing and exploring information.
You can search for users of specific type (associations), trainings or events, locations.
As users can have translatable descriptions, just like trainings and locations,
a full text search will be available. Locations will use geospatial sorting, and because
trainings happen at locations, they can also be sorted by distance.

Based on this, the navigation's purpose is to load the correct filters and search results.
The user is given the following choices:

* `[Entities]`
* `[Trainings]`
* `[Events]`
* `[Locations]`

The search results of trainings, events and locations can also be displayed in a map view.
Views should remind users of a hotel booking website, where the map is on the right side,
and the search results are on the left side. The map should be interactive, and panning around
should update the search results. The search results should be sorted by distance to the center
of the map. The map should be centered on the user's current location, if available.

The user can also choose to view the search results in a list view, which is the default.
Here, the user can use filter criteria like the maximum distance, search for specific search terms,
and filtering trainings by an entity should also be possible.

This also means, if a user found an interesting entity, in some kind of detail page,
the user should be able to query all trainings and events of this entity,
opening the search results page with the entity preselected.

Likewise, if a user found an interesting location, the user could see if there are any
trainings or events happening there, and if so, when.

A training is usually conducted by an entity, and the user should be able to see
details about the entity, and could also check for other trainings of this entity.
The location obviously also plays a role, and the user should be able to see
details about the location, ideally a map screen is shown with the location.

There isn't much difference between a training and an event. I would assume, an event is
something special that doesn't happen monthly or even more frequently.

This allows for a training calendar or event calendar view, and because events
are usually more interesting for users living far away, training calendars are usually
region- or entity-specific, and event calendars are usually country- or world-wide.

We shouldn't forget to mention that a user that is logged in can create their own trainings,
events and locations. They can even edit locations, add pictures or description to them,
add comments, and so on. As entities can also have pages, these pages can act as
a personal website for the entity with subpages and subsections.

With the search results being able to be presented as a list, a map or a calendar,
and because the data needed to compile the calendar view is slightly different,
the calendar view should be a separate page, and the user should be able to switch
between the list, map and calendar view.

```text
----------------------
|   Filters          |
----------------------
|  Search Terms:     |
|   [_____________]  |
----------------------
|  Type:             |
|   [ ] Entities     |
|   [x] Trainings    |
|   [ ] Events       |
|   [ ] Locations    |
----------------------
|  View:             |
|   [x] Map          |
|   [ ] List         |
|   [ ] Calendar     |
----------------------
|  Location:         |
|   [x] GPS location |
|   [ ] Selected     |
|   [ ] Address      |
|       [_________]  |
----------------------
|  Maximum Distance: |
|   [   5 km      ]  |
----------------------
|  Organiser:        |
|   [_____________]  |
----------------------
|  Type of training: |
|   [x] parkour      |
|   [ ] parkour-jam  |
|   [ ] meeting      |
|   [ ] show         |
|   [ ] competition  |
|   [ ] slackline    |
|   [ ] tour         |
----------------------
|  Weekday       [v] |
----------------------
|  Time of day   [v] |
----------------------
|   [ Apply ]        |
----------------------
```