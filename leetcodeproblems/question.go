package leetcodeproblems

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type QuestionDatabase struct {
	Questions               map[string][]string
	QuestionsFinishSummary  map[string][]*QuestionFinishSummary
	QuestionsFinishDailyLog []*QuestionFinishDailyLog
}

type QuestionFinishSummary struct {
	Name           string
	QuestionType   string
	FinishTimes    uint64
	LastFinishTime int64
}

type QuestionFinishDailyLog struct {
	Name       string
	FinishTime int64
}

func NewQuestionDatabase() QuestionDatabase {
	var (
		questionDatabase = QuestionDatabase{}
	)

	var (
		err      error
		question = make(map[string][]string)
	)
	question, err = GetAllQuestions()
	if err != nil {
		panic(err)
	}
	questionDatabase.Questions = question

	var (
		questionFinishDailyLogs []*QuestionFinishDailyLog
	)
	questionFinishDailyLogs, err = GetAllQuestionFinishDailyLogs()
	if err != nil {
		panic(err)
	}
	questionDatabase.QuestionsFinishDailyLog = questionFinishDailyLogs

	var (
		questionFinishSummary map[string][]*QuestionFinishSummary
	)
	questionFinishSummary, err = GetAllQuestionsFinishSummary(question, questionFinishDailyLogs)
	if err != nil {
		panic(err)
	}
	questionDatabase.QuestionsFinishSummary = questionFinishSummary

	return questionDatabase
}

func GetAllQuestions() (map[string][]string, error) {
	var (
		questionFile *os.File
		err          error
	)
	questionFile, err = os.Open("data/question.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		err = questionFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	var (
		byteValue []byte
	)
	byteValue, err = io.ReadAll(questionFile)
	if err != nil {
		return nil, err
	}

	var (
		questions = make(map[string][]string)
	)
	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func GetAllQuestionFinishDailyLogs() ([]*QuestionFinishDailyLog, error) {
	var (
		questionFinishDailyLog []*QuestionFinishDailyLog
	)
	dl, err := os.Open("data/daily_log.log")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = dl.Close()
		if err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(dl)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			panic(fmt.Errorf("illegal row: %v", line))
		}
		finishTime, _ := strconv.ParseInt(fields[1], 10, 64)
		questionFinishDailyLog = append(questionFinishDailyLog, &QuestionFinishDailyLog{
			Name:       fields[0],
			FinishTime: finishTime,
		})
	}
	return questionFinishDailyLog, nil
}

func GetAllQuestionsFinishSummary(questions map[string][]string, questionsFinishDailyLog []*QuestionFinishDailyLog) (map[string][]*QuestionFinishSummary, error) {
	var (
		questionsFinishSummary = make(map[string][]*QuestionFinishSummary)
	)
	for t, questionNames := range questions {
		for i := 0; i < len(questionNames); i++ {
			questionsFinishSummary[t] = append(questionsFinishSummary[t], &QuestionFinishSummary{
				Name:           questionNames[i],
				QuestionType:   t,
				FinishTimes:    0,
				LastFinishTime: 0,
			})
		}
	}

	for i := 0; i < len(questionsFinishDailyLog); i++ {
		var (
			name = questionsFinishDailyLog[i].Name
		)
		for _, summaries := range questionsFinishSummary {
			for j := 0; j < len(summaries); j++ {
				if name == summaries[j].Name {
					if summaries[j].LastFinishTime < questionsFinishDailyLog[i].FinishTime {
						summaries[j].LastFinishTime = questionsFinishDailyLog[i].FinishTime
					}
					summaries[j].FinishTimes++
				}
			}
		}
	}

	return questionsFinishSummary, nil
}

func (q *QuestionDatabase) Random() error {
	var (
		list  []string
		count int
	)
	for t, summaries := range q.QuestionsFinishSummary {
		if t == "Must" {
			for i := 0; i < len(summaries); i++ {
				list = append(list, summaries[i].Name)
			}
			continue
		}
		sort.Slice(summaries, func(i, j int) bool {
			return summaries[i].LastFinishTime < summaries[j].LastFinishTime
		})
		if time.Now().Unix()-summaries[0].LastFinishTime >= 86400 && count < 5 {
			list = append(list, summaries[0].Name)
			count++
		}
	}
	for i := 0; i < len(list); i++ {
		fmt.Println(list[i])
		_ = exec.Command("open", "https://leetcode.cn/problems/"+list[i]).Start()
	}

	var (
		err      error
		dailyLog *os.File
	)
	dailyLog, err = os.OpenFile("data/daily_log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = dailyLog.Close()
		if err != nil {
			panic(err)
		}
	}()
	for i := 0; i < len(list); i++ {
		_, err = dailyLog.WriteString(list[i] + ",")
		if err != nil {
			panic(err)
		}
		_, err = dailyLog.WriteString(strconv.Itoa(int(time.Now().Unix())) + ",")
		if err != nil {
			panic(err)
		}
		_, err = dailyLog.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (q *QuestionDatabase) PrintAllQuestionsFinishSummary() error {
	for _, questions := range q.QuestionsFinishSummary {
		for i := 0; i < len(questions); i++ {
			finishSummary := questions[i]
			fmt.Printf("%v,%v,%v,%v\n", finishSummary.Name, finishSummary.QuestionType, finishSummary.FinishTimes, time.Unix(finishSummary.LastFinishTime, 0).Format("2006-01-02 15:04:05"))
		}
	}
	return nil
}
