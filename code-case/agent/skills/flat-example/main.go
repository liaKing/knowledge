package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SkillIndex struct {
	Name        string
	Description string
	Dir         string
}

type WeatherResult struct {
	City        string `json:"city"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
	Wind        string `json:"wind"`
	Advice      string `json:"advice"`
}

func main() {
	userRequest := "北京今天适合出门吗？"
	fmt.Println("用户请求:", userRequest)

	skills, err := discoverFlatSkills(".")
	if err != nil {
		panic(err)
	}
	printSkillIndex(skills)

	matched, ok := routeFlatSkill(userRequest, skills)
	if !ok {
		fmt.Println("没有命中 Skill。")
		return
	}
	fmt.Printf("\n命中 Skill: %s\n", matched.Name)

	instructions, err := os.ReadFile(filepath.Join(matched.Dir, "SKILL.md"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("按渐进式披露加载完整 SKILL.md: %s (%d bytes)\n", filepath.Join(matched.Dir, "SKILL.md"), len(instructions))

	switch matched.Name {
	case "weather-advice":
		result, err := runWeatherSkill(matched.Dir, extractCity(userRequest))
		if err != nil {
			panic(err)
		}
		fmt.Println("脚本返回:", mustJSON(result))
		fmt.Printf("最终回答: %s今天%s，%s，%s。%s\n", result.City, result.Weather, result.Temperature, result.Wind, result.Advice)
	case "git-commit-message":
		fmt.Println("最终回答: feat(agent): add skill examples")
	default:
		fmt.Println("命中 Skill，但 demo 未实现对应 executor。")
	}
}

func discoverFlatSkills(root string) ([]SkillIndex, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var skills []SkillIndex
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dir := filepath.Join(root, entry.Name())
		content, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
		if err != nil {
			continue
		}

		name, description := parseFrontmatter(string(content))
		if name == "" {
			continue
		}
		skills = append(skills, SkillIndex{
			Name:        name,
			Description: description,
			Dir:         dir,
		})
	}
	return skills, nil
}

func routeFlatSkill(request string, skills []SkillIndex) (SkillIndex, bool) {
	for _, skill := range skills {
		if skill.Name == "weather-advice" && containsAny(request, "天气", "温度", "出门", "降雨", "风") {
			return skill, true
		}
		if skill.Name == "git-commit-message" && containsAny(request, "commit", "提交", "变更") {
			return skill, true
		}
	}
	return SkillIndex{}, false
}

func runWeatherSkill(skillDir string, city string) (WeatherResult, error) {
	cmd := exec.Command("go", "run", "scripts/fixed_weather.go", "--city", city)
	cmd.Dir = skillDir

	output, err := cmd.Output()
	if err != nil {
		return WeatherResult{}, err
	}

	var result WeatherResult
	if err := json.Unmarshal(output, &result); err != nil {
		return WeatherResult{}, err
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

func printSkillIndex(skills []SkillIndex) {
	fmt.Println("\n第一阶段：只暴露 Skill 索引")
	for _, skill := range skills {
		fmt.Printf("- %s: %s\n", skill.Name, skill.Description)
	}
}

func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func extractCity(text string) string {
	for _, city := range []string{"北京", "上海", "深圳", "广州", "杭州"} {
		if strings.Contains(text, city) {
			return city
		}
	}
	return "北京"
}

func mustJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
