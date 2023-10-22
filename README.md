# One Hundred Year Decay

One Hundred Year Decay is a project done for Week 4 of The Musical Web, a class at the School for Poetic Computation taught by Chloe Alexandra Thompson & Tommy Martinez.

One Hundred Year Decay is heavily inspired by William Basinki's Disintegration Loops, and is an attempt to recreate them on a much more glacial time scale.

A single key/value pair is decremented outside of the browser, and when you listen to the decay, this value is retrieved and used to configure various settings for effects in Tone.JS.

### Tech

This project was mostly an attempt to learn Go, after doing a large portion of the book, [Let's Go](https://lets-go.alexedwards.net/). It uses the Go CDK bindings to provision infrastructure, a Go lambda to slowly decrement our value, and an unbelievably simple web server using Gin.

Go is fun language! It is absurdly overpowered for this task, but I think I'll continue to use it for backend projects.

You can run this in your own AWS environment if you'd like, by calling `cdk deploy` while in the `infra` directory. 