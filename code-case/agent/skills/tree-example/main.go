package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SkillNode struct {
	Name        string
	Description string
	Dir         string
}

type BeachWeatherResult struct {
	Place       string `json:"place"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
	Wind        string `json:"wind"`
	UV          string `json:"uv"`
	Risk        string `json:"risk"`
	Advice      string `json:"advice"`
}

func main() {
	userRequest := "周末去青岛海边玩，天气适合吗？"
	fmt.Println("用户请求:", userRequest)

	root, err := loadSkillNode("weather")
	if err != nil {
		panic(err)
	}

	fmt.Println("\n第一阶段：只暴露顶层 Skill 索引")
	fmt.Printf("- %s: %s\n", root.Name, root.Description)

	fmt.Printf("\n命中顶层 Skill: %s\n", root.Name)
	printLoadedSkill(root)

	travel, err := loadSkillNode(filepath.Join(root.Dir, "travel"))
	if err != nil {
		panic(err)
	}
	fmt.Println("根据 weather 路由规则，命中子 Skill: travel")
	printLoadedSkill(travel)

	beach, err := loadSkillNode(filepath.Join(travel.Dir, "beach"))
	if err != nil {
		panic(err)
	}
	fmt.Println("根据 travel-weather 路由规则，命中叶子 Skill: beach")
	printLoadedSkill(beach)

	result, err := runBeachWeatherSkill(beach.Dir, extractPlace(userRequest))
	if err != nil {
		panic(err)
	}

	fmt.Println("脚本返回:", mustJSON(result))
	fmt.Printf("最终回答: %s今天%s，约%s，%s，紫外线%s。%s，%s\n",
		result.Place,
		result.Weather,
		result.Temperature,
		result.Wind,
		result.UV,
		result.Risk,
		result.Advice,
	)
}

func loadSkillNode(dir string) (SkillNode, error) {
	content, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	if err != nil {
		return SkillNode{}, err
	}

	name, description := parseFrontmatter(string(content))
	if name == "" {
		return SkillNode{}, fmt.Errorf("missing skill name in %s", dir)
	}

	return SkillNode{
		Name:        name,
		Description: description,
		Dir:         dir,
	}, nil
}

func printLoadedSkill(skill SkillNode) {
	content, err := os.ReadFile(filepath.Join(skill.Dir, "SKILL.md"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("按渐进式披露加载完整 SKILL.md: %s (%d bytes)\n", filepath.Join(skill.Dir, "SKILL.md"), len(content))
}

func runBeachWeatherSkill(skillDir string, place string) (BeachWeatherResult, error) {
	cmd := exec.Command("go", "run", "scripts/fixed_beach_weather.go", "--place", place)
	cmd.Dir = skillDir

	output, err := cmd.Output()
	if err != nil {
		return BeachWeatherResult{}, err
	}

	var result BeachWeatherResult
	if err := json.Unmarshal(output, &result); err != nil {
		return BeachWeatherResult{}, err
	}
	return result, nil
}

func parseFrontmatter(markdown string) (string, string) {
	if !strings.HasPrefix(markdown, "---") {
		return "", ""
	}

	parts := strings.SplitN(markdown, "---", 3)
	if len(parts) < 3 {
		return "", ""
	}

	var name string
	var description string
	for _, line := range strings.Split(parts[1], "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "name:"):
			name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		case strings.HasPrefix(line, "description:"):
			description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
	}
	return name, description
}

func extractPlace(text string) string {
	for _, place := range []string{"青岛海边", "三亚海边", "厦门海边", "舟山海边"} {
		if strings.Contains(text, strings.TrimSuffix(place, "海边")) {
			return place
		}
	}
	return "青岛海边"
}

func mustJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
