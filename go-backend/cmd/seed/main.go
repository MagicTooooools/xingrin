package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/xingrin/go-backend/internal/config"
	"github.com/xingrin/go-backend/internal/database"
	"github.com/xingrin/go-backend/internal/model"
	"gorm.io/gorm"
)

var (
	clear       = flag.Bool("clear", false, "Clear existing data before generating")
	orgCount    = flag.Int("orgs", 20, "Number of organizations to generate")
	targetCount = flag.Int("targets", 100, "Number of targets to generate")
)

func main() {
	flag.Parse()

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
	fmt.Printf("   Targets: %d\n", *targetCount)
	fmt.Println()

	if *clear {
		fmt.Println("🗑️  Clearing existing data...")
		if err := clearData(db); err != nil {
			fmt.Printf("❌ Failed to clear data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("   ✓ Data cleared\n")
	}

	// Generate data
	orgIDs, err := createOrganizations(db, *orgCount)
	if err != nil {
		fmt.Printf("❌ Failed to create organizations: %v\n", err)
		os.Exit(1)
	}

	targetIDs, err := createTargets(db, *targetCount)
	if err != nil {
		fmt.Printf("❌ Failed to create targets: %v\n", err)
		os.Exit(1)
	}

	// Link targets to organizations
	if err := linkTargetsToOrganizations(db, targetIDs, orgIDs); err != nil {
		fmt.Printf("❌ Failed to link targets to organizations: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Test data generation completed!")
}

func clearData(db *gorm.DB) error {
	// Delete in order to respect foreign key constraints
	tables := []string{
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

	// Each target belongs to 1-3 organizations
	linkCount := 0
	for _, targetID := range targetIDs {
		numOrgs := rand.Intn(3) + 1 // 1-3 organizations
		selectedOrgs := make(map[int]bool)

		for i := 0; i < numOrgs; i++ {
			orgID := orgIDs[rand.Intn(len(orgIDs))]
			if selectedOrgs[orgID] {
				continue
			}
			selectedOrgs[orgID] = true

			err := db.Exec(
				"INSERT INTO organization_target (organization_id, target_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				orgID, targetID,
			).Error
			if err != nil {
				return err
			}
			linkCount++
		}
	}

	fmt.Printf("   ✓ Created %d target-organization links\n", linkCount)
	return nil
}
