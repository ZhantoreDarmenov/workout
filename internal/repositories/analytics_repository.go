package repositories

import (
	"context"
	"database/sql"

	"workout/internal/models"
)

// AnalyticsRepository provides methods to gather statistics for trainers.
type AnalyticsRepository struct {
	DB *sql.DB
}

// TrainerAnalytics computes various metrics for a trainer.
func (r *AnalyticsRepository) TrainerAnalytics(ctx context.Context, trainerID int) (models.TrainerAnalytics, error) {
	var res models.TrainerAnalytics

	err := r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM workout_programs WHERE trainer_id = ?`, trainerID).Scan(&res.ProgramCount)
	if err != nil {
		return res, err
	}

	err = r.DB.QueryRowContext(ctx, `SELECT COUNT(DISTINCT p.client_id)
        FROM progress p
        JOIN days d ON p.day_id = d.id
        JOIN workout_programs wp ON d.work_out_program_id = wp.id
        WHERE wp.trainer_id = ?`, trainerID).Scan(&res.ClientCount)
	if err != nil {
		return res, err
	}

	rows, err := r.DB.QueryContext(ctx, `SELECT p.client_id,
        SUM(CASE WHEN p.completed IS NOT NULL THEN 1 ELSE 0 END) AS completed_days,
        COUNT(d.id) AS total_days
        FROM progress p
        JOIN days d ON p.day_id = d.id
        JOIN workout_programs wp ON d.work_out_program_id = wp.id
        WHERE wp.trainer_id = ?
        GROUP BY p.client_id`, trainerID)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	var progressSum float64
	for rows.Next() {
		var cp models.ClientProgress
		if err := rows.Scan(&cp.ClientID, &cp.CompletedDays, &cp.TotalDays); err != nil {
			return res, err
		}
		if cp.TotalDays > 0 {
			cp.Progress = float64(cp.CompletedDays) / float64(cp.TotalDays)
		}
		if cp.Progress == 1 {
			res.CompletedClients++
		}
		progressSum += cp.Progress
		res.ClientsProgress = append(res.ClientsProgress, cp)
	}
	if err := rows.Err(); err != nil {
		return res, err
	}

	if len(res.ClientsProgress) > 0 {
		res.AverageProgress = progressSum / float64(len(res.ClientsProgress))
	}

	return res, nil
}
