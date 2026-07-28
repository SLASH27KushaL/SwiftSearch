# SwiftSearch 🔎

Hey there! Welcome to **SwiftSearch**—a highly scalable, distributed search engine I built from scratch to dive deep into microservices, algorithmic ranking, and complex system design. 

Think of it as a custom, lightweight Google clone. It actively crawls the web, parses HTML, builds an inverted index, and calculates link authority in the background. I engineered this to tackle real-world backend challenges like managing distributed queues, bypassing bot-blockers, and writing optimized mathematical sorting expressions. 

## 🔗 Live Demo
You can check out the repository and upcoming live deployment here: **[SwiftSearch](https://github.com/SLASH27KushaL/SwiftSearch)**

---

## ✨ What it does

*   **Custom Web Spider:** A distributed crawler that respects `robots.txt`, enforces politeness delays, and uses a Bloom Filter to avoid redundant scraping. It even spoofs User-Agents to fetch secure HTML gracefully.
*   **Algorithmic Ranking:** Results aren't just keyword matches. The engine utilizes a mathematical fusion model to sort search results, balancing keyword relevance against domain authority: **Score = (α * TF) + (β * PR)**.
*   **Inverted Indexing:** Tokenizes documents, removes stopwords, and calculates Term Frequency (TF) for lightning-fast keyword lookups.
*   **PageRank Engine:** A background job that traverses the inter-document link graph to assign objective authority scores (PR) to URLs.
*   **Real-Time API:** A lightning-fast Go backend that intercepts queries and fetches perfectly ranked JSON payloads, utilizing Redis to serve cached hits in sub-milliseconds.

---

## 🛠️ The Tech Stack

I leaned into a highly decoupled microservices architecture for this one. Here's what's under the hood:

*   **Backend & Workers:** Go (Golang) & Gin Web Framework. The crawler, indexer, and search API are all written in Go for maximum concurrency and performance.
*   **Frontend Framework:** Next.js & React, styled cleanly to render dynamic, real-time search results without caching stale data.
*   **Databases:** MongoDB (handles all document storage, inverted index mapping, and PageRank scores).
*   **Caching & Queues:** Redis (manages the distributed frontier queue for the spider and aggressively caches Search API responses).
*   **Infrastructure:** Docker & Docker Compose. The entire system is containerized so all services can be spun up seamlessly.

---

## 💻 Local Development Setup

Want to spin this entire distributed system up on your own machine? Here is how to get it running.

### 1. Prerequisites
Make sure you have **Docker** and **Docker Compose** installed on your machine. You will also need **Node.js** for the frontend.

### 2. Clone and Install
```bash
git clone [https://github.com/SLASH27KushaL/SwiftSearch.git](https://github.com/SLASH27KushaL/SwiftSearch.git)
cd SwiftSearch
