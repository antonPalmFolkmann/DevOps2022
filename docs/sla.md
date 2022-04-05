# Service Level Agreement

## Introduction

Minitwit is an educational service provided for free.
This is an educational document written to better understand SLAs.
Given the groups work schedule we are best equipped to fix irregularities with the service on Tuesdays and Sundays.

## General Terms

- Monthly Uptime: The number of hours in the weekday the service is down / the number of hours in the weekday the service is running
- Running: When all endpoints are able to return 2xx status code
- Down: When an endpoint returns a non-2xx status code
- Back-off Period: In case of a failed response the client is responsible for backing off. The minimum back-off period is 10 seconds and the client should wait longer and longer between each request

### Limitations

- We reserve the ability to expend up to 72 hours to debug/solve an issue
- You are only eligble to claim credits once a month

## Monthly Uptime

| Uptime | Credit |
|:------:|:------:|
| <= 80% | A beautiful folded paperplane |
| <= 75% | A handwritten apology note signed by the entire team |
| <= 50% | A homemade cake |

## Call metrics

We would love to hear from you, you can reach us at:

- Github Issues
- Our ITU Mails
