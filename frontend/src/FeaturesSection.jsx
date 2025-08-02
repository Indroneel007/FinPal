import React from 'react';
import { Box, Grid, Paper, Typography, useTheme } from '@mui/material';
import SavingsIcon from '@mui/icons-material/Savings';
import GroupIcon from '@mui/icons-material/Group';
import TimelineIcon from '@mui/icons-material/Timeline';
import SecurityIcon from '@mui/icons-material/Security';
import AutoGraphIcon from '@mui/icons-material/AutoGraph';
import CreditScoreIcon from '@mui/icons-material/CreditScore';

const features = [
  {
    icon: <SavingsIcon fontSize="large" color="primary" />, title: 'Smart Budgeting', desc: 'Track income, expenses, and savings goals with ease.'
  },
  {
    icon: <GroupIcon fontSize="large" color="secondary" />, title: 'Group Transactions', desc: 'Manage shared expenses and group payments seamlessly.'
  },
  {
    icon: <TimelineIcon fontSize="large" color="action" />, title: 'Transaction History', desc: 'Visualize spending patterns and transaction history.'
  },
  {
    icon: <SecurityIcon fontSize="large" color="success" />, title: 'Secure & Private', desc: 'Your financial data is encrypted and protected.'
  },
  {
    icon: <AutoGraphIcon fontSize="large" color="warning" />, title: 'Analytics & Insights', desc: 'Get personalized financial insights and AI-powered tips.'
  },
  {
    icon: <CreditScoreIcon fontSize="large" color="error" />, title: 'Multi-Account Support', desc: 'Handle multiple accounts and currencies with one app.'
  }
];

export default function FeaturesSection() {
  const theme = useTheme();
  return (
    <Box sx={{ width: '100%', bgcolor: 'transparent', py: 8 }}>
      <Typography variant="h4" align="center" fontWeight={700} color={theme.palette.primary.light} mb={4}>
        Why Choose FinPal?
      </Typography>
      <Grid container spacing={4} justifyContent="center">
        {features.map((f, i) => (
          <Grid item xs={12} sm={6} md={4} key={f.title}>
            <Paper elevation={4} sx={{ p: 3, borderRadius: 3, textAlign: 'center', transition: '0.2s', '&:hover': { boxShadow: 10, transform: 'scale(1.04)' }, minHeight: 200 }}>
              <Box mb={2} display="flex" justifyContent="center" alignItems="center">
                {f.icon}
              </Box>
              <Typography variant="h6" fontWeight={600} gutterBottom>{f.title}</Typography>
              <Typography variant="body2" color="text.secondary">{f.desc}</Typography>
            </Paper>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
}
