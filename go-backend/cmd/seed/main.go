package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/xingrin/go-backend/internal/config"
	"github.com/xingrin/go-backend/internal/database"
	"github.com/xingrin/go-backend/internal/model"
	"gorm.io/gorm"
)

var (
	clear    = flag.Bool("clear", false, "Clear existing data before generating")
	orgCount = flag.Int("orgs", 20, "Number of organizations to generate")
	// targetCount = orgCount * 20, websiteCount = targetCount * 20
)

func main() {
	flag.Parse()

	// Calculate counts based on org count
	// Each org has 20 targets, each target has 20 websites
	targetsPerOrg := 20
	websitesPerTarget := 20
	targetCount := *orgCount * targetsPerOrg

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Starting test data generation...")
	fmt.Printf("   Organizations: %d\n", *orgCount)
	fmt.Printf("   Targets: %d (%d per org)\n", targetCount, targetsPerOrg)
	fmt.Printf("   Websites: %d (%d per target)\n", targetCount*websitesPerTarget, websitesPerTarget)
	fmt.Println()

	if *clear {
		fmt.Println("🗑️  Clearing existing data...")
		if err := clearData(db); err != nil {
			fmt.Printf("❌ Failed to clear data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("   ✓ Data cleared")
	}

	// Generate data
	orgIDs, err := createOrganizations(db, *orgCount)
	if err != nil {
		fmt.Printf("❌ Failed to create organizations: %v\n", err)
		os.Exit(1)
	}

	targetIDs, err := createTargets(db, targetCount)
	if err != nil {
		fmt.Printf("❌ Failed to create targets: %v\n", err)
		os.Exit(1)
	}

	// Link targets to organizations (20 per org)
	if err := linkTargetsToOrganizations(db, targetIDs, orgIDs); err != nil {
		fmt.Printf("❌ Failed to link targets to organizations: %v\n", err)
		os.Exit(1)
	}

	// Create websites for targets (20 per target)
	if err := createWebsites(db, targetIDs, websitesPerTarget); err != nil {
		fmt.Printf("❌ Failed to create websites: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Test data generation completed!")
}

func clearData(db *gorm.DB) error {
	// Delete in order to respect foreign key constraints
	tables := []string{
		"website",
		"organization_target",
		"target",
		"organization",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return fmt.Errorf("failed to clear %s: %w", table, err)
		}
	}
	return nil
}

func createOrganizations(db *gorm.DB, count int) ([]int, error) {
	fmt.Printf("🏢 Creating %d organizations...\n", count)

	orgNames := []string{
		"Acme Corporation", "TechStart Labs", "Global Finance", "HealthCare Plus",
		"E-Commerce Platform", "Smart City Systems", "Educational Tech", "Green Energy",
		"CyberSec Defense", "CloudNative Systems", "DataFlow Analytics", "MobileFirst Tech",
		"Quantum Research", "Autonomous Vehicles", "Biotech Innovations", "Space Technology",
		"AI Research Lab", "Blockchain Solutions", "IoT Platform", "DevOps Enterprise",
		"Security Operations", "Data Science Hub", "Machine Learning Co", "Network Solutions",
		"Infrastructure Corp", "Platform Services", "Digital Transformation", "Innovation Hub",
		"Tech Consulting", "Software Factory",
	}

	divisions := []string{
		"Global", "Asia Pacific", "EMEA", "Americas", "R&D", "Cloud Services",
		"Security Team", "Innovation Lab", "Enterprise", "Consumer Products",
	}

	descriptions := []string{
		"A leading technology company specializing in enterprise software solutions and cloud computing services.",
		"Innovative research lab focused on artificial intelligence and machine learning applications.",
		"Global financial services provider offering digital banking and payment solutions.",
		"Healthcare technology company developing electronic health records and telemedicine platforms.",
		"E-commerce platform serving millions of customers with B2B and B2C solutions.",
		"Smart city infrastructure provider specializing in IoT and urban management systems.",
		"Educational technology company providing online learning platforms and courses.",
		"Renewable energy management company focused on solar and wind power optimization.",
		"Cybersecurity firm offering penetration testing and security consulting services.",
		"Cloud-native systems developer specializing in Kubernetes and microservices.",
	}

	var ids []int
	suffix := rand.Intn(9000) + 1000

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%s - %s (%d-%d)",
			orgNames[i%len(orgNames)],
			divisions[rand.Intn(len(divisions))],
			suffix, i)

		desc := descriptions[rand.Intn(len(descriptions))]

		// Random created_at within last year
		daysAgo := rand.Intn(365)
		createdAt := time.Now().AddDate(0, 0, -daysAgo)

		org := &model.Organization{
			Name:        name,
			Description: desc,
			CreatedAt:   createdAt,
		}

		if err := db.Create(org).Error; err != nil {
			return nil, err
		}
		ids = append(ids, org.ID)
	}

	fmt.Printf("   ✓ Created %d organizations\n", len(ids))
	return ids, nil
}

func createTargets(db *gorm.DB, count int) ([]int, error) {
	fmt.Printf("🎯 Creating %d targets...\n", count)

	// Domain templates
	envs := []string{"prod", "staging", "dev", "test", "api", "app", "www", "admin", "portal", "dashboard"}
	companies := []string{"acme", "techstart", "globalfinance", "healthcare", "ecommerce", "smartcity", "cybersec", "cloudnative", "dataflow", "mobilefirst"}
	tlds := []string{".com", ".io", ".net", ".org", ".dev", ".app", ".cloud", ".tech"}

	var ids []int
	suffix := rand.Intn(9000) + 1000
	usedNames := make(map[string]bool)

	// Generate domains (70%)
	domainCount := count * 70 / 100
	for i := 0; i < domainCount; i++ {
		var name string
		for {
			env := envs[rand.Intn(len(envs))]
			company := companies[rand.Intn(len(companies))]
			tld := tlds[rand.Intn(len(tlds))]
			name = fmt.Sprintf("%s.%s-%d%s", env, company, suffix+i, tld)
			if !usedNames[name] {
				usedNames[name] = true
				break
			}
		}

		target := &model.Target{
			Name:      name,
			Type:      "domain",
			CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(365)),
		}

		if err := db.Create(target).Error; err != nil {
			return nil, err
		}
		ids = append(ids, target.ID)
	}

	// Generate IPs (20%)
	ipCount := count * 20 / 100
	for i := 0; i < ipCount; i++ {
		name := fmt.Sprintf("%d.%d.%d.%d",
			rand.Intn(223)+1,
			rand.Intn(256),
			rand.Intn(256),
			rand.Intn(254)+1)

		if usedNames[name] {
			continue
		}
		usedNames[name] = true

		target := &model.Target{
			Name:      name,
			Type:      "ip",
			CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(365)),
		}

		if err := db.Create(target).Error; err != nil {
			return nil, err
		}
		ids = append(ids, target.ID)
	}

	// Generate CIDRs (10%)
	cidrCount := count * 10 / 100
	for i := 0; i < cidrCount; i++ {
		masks := []int{8, 16, 24}
		mask := masks[rand.Intn(len(masks))]
		name := fmt.Sprintf("%d.%d.%d.0/%d",
			rand.Intn(223)+1,
			rand.Intn(256),
			rand.Intn(256),
			mask)

		if usedNames[name] {
			continue
		}
		usedNames[name] = true

		target := &model.Target{
			Name:      name,
			Type:      "cidr",
			CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(365)),
		}

		if err := db.Create(target).Error; err != nil {
			return nil, err
		}
		ids = append(ids, target.ID)
	}

	fmt.Printf("   ✓ Created %d targets (domains: %d, IPs: %d, CIDRs: %d)\n",
		len(ids), domainCount, ipCount, cidrCount)
	return ids, nil
}

func linkTargetsToOrganizations(db *gorm.DB, targetIDs, orgIDs []int) error {
	fmt.Println("🔗 Linking targets to organizations...")

	if len(orgIDs) == 0 || len(targetIDs) == 0 {
		return nil
	}

	// Each organization gets exactly 20 targets (evenly distributed)
	targetsPerOrg := len(targetIDs) / len(orgIDs)
	linkCount := 0

	for orgIdx, orgID := range orgIDs {
		startIdx := orgIdx * targetsPerOrg
		endIdx := startIdx + targetsPerOrg
		if orgIdx == len(orgIDs)-1 {
			endIdx = len(targetIDs) // Last org gets remaining targets
		}

		for i := startIdx; i < endIdx; i++ {
			err := db.Exec(
				"INSERT INTO organization_target (organization_id, target_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				orgID, targetIDs[i],
			).Error
			if err != nil {
				return err
			}
			linkCount++
		}
	}

	fmt.Printf("   ✓ Created %d target-organization links (%d per org)\n", linkCount, targetsPerOrg)
	return nil
}

func createWebsites(db *gorm.DB, targetIDs []int, websitesPerTarget int) error {
	totalCount := len(targetIDs) * websitesPerTarget
	fmt.Printf("🌐 Creating %d websites (%d per target)...\n", totalCount, websitesPerTarget)

	if len(targetIDs) == 0 {
		return nil
	}

	// Website data templates
	protocols := []string{"https://", "http://"}
	subdomains := []string{"www", "api", "app", "admin", "portal", "dashboard", "dev", "staging", "test", "cdn", "static", "assets", "mail", "blog", "docs", "support", "auth", "login", "shop", "store"}
	paths := []string{"", "/", "/api", "/v1", "/v2", "/login", "/dashboard", "/admin", "/app", "/docs"}
	ports := []string{"", ":8080", ":8443", ":3000", ":443"}

	titles := []string{
		"Welcome - Dashboard", "Admin Panel", "API Documentation", "Login Portal",
		"Home Page", "User Dashboard", "Settings", "Analytics", "Reports",
		"Management Console", "Control Panel", "Service Status", "Developer Portal",
	}

	webservers := []string{
		"nginx/1.24.0", "Apache/2.4.57", "cloudflare", "Microsoft-IIS/10.0",
		"nginx", "Apache", "LiteSpeed", "Caddy", "Traefik",
	}

	contentTypes := []string{
		"text/html; charset=utf-8", "text/html", "application/json",
		"text/html; charset=UTF-8", "application/xhtml+xml",
	}

	techStacks := [][]string{
		{"nginx", "PHP", "MySQL"},
		{"Apache", "Python", "PostgreSQL"},
		{"nginx", "Node.js", "MongoDB"},
		{"cloudflare", "React", "GraphQL"},
		{"nginx", "Vue.js", "Redis"},
		{"Apache", "Java", "Oracle"},
		{"nginx", "Go", "PostgreSQL"},
		{"cloudflare", "Next.js", "Vercel"},
		{"nginx", "Django", "PostgreSQL"},
		{"Apache", "Laravel", "MySQL"},
		{"nginx", "Express", "MongoDB"},
		{"cloudflare", "Angular", "Firebase"},
	}

	statusCodes := []int{200, 200, 200, 200, 200, 301, 302, 403, 404, 500}

	createdCount := 0

	// Each target gets exactly websitesPerTarget websites
	for targetIdx, targetID := range targetIDs {
		for i := 0; i < websitesPerTarget; i++ {
			// Generate deterministic but varied URL
			protocol := protocols[i%len(protocols)]
			subdomain := subdomains[i%len(subdomains)]
			port := ports[i%len(ports)]
			path := paths[i%len(paths)]

			domain := fmt.Sprintf("target%d-site%d.com", targetIdx, i)
			url := fmt.Sprintf("%s%s.%s%s%s", protocol, subdomain, domain, port, path)

			// Extract host from URL
			host := extractHost(url)

			// Deterministic but varied values
			statusCode := statusCodes[i%len(statusCodes)]
			contentLength := 1000 + (i * 100)
			tech := pq.StringArray(techStacks[i%len(techStacks)])
			vhost := i%5 == 0 // 20% are vhost

			website := &model.Website{
				TargetID:      targetID,
				URL:           url,
				Host:          host,
				Title:         titles[i%len(titles)],
				StatusCode:    &statusCode,
				ContentLength: &contentLength,
				ContentType:   contentTypes[i%len(contentTypes)],
				Webserver:     webservers[i%len(webservers)],
				Tech:          tech,
				Vhost:         &vhost,
				Location:      "",
				CreatedAt:     time.Now().AddDate(0, 0, -i),
			}

			if err := db.Create(website).Error; err != nil {
				continue
			}
			createdCount++
		}
	}

	fmt.Printf("   ✓ Created %d websites\n", createdCount)
	return nil
}

func extractHost(rawURL string) string {
	// Simple host extraction
	url := rawURL
	// Remove protocol
	if idx := len("https://"); len(url) > idx && url[:idx] == "https://" {
		url = url[idx:]
	} else if idx := len("http://"); len(url) > idx && url[:idx] == "http://" {
		url = url[idx:]
	}
	// Remove path
	if idx := findIndex(url, "/"); idx != -1 {
		url = url[:idx]
	}
	return url
}

func findIndex(s string, substr string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == substr[0] {
			return i
		}
	}
	return -1
}
